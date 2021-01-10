---
title: "My Go proverbs"
publishdate: 2020-11-17
resources:
    - name: header
    - src: featured.jpg
categories:
    - Golang
tags:
    - system design
    - tests
---

I prepared a list of my [Go proverbs](https://go-proverbs.github.io/). This list contains 7 points I try to follow while working on Go code. The list is based on my personal experience and you may not agree with some of them. That's fine! Just let us know in the comments section about it or just share with us with your own suggestions.

## Focus on what your API does, not what it is

To understand what I mean here we have to think about the API from two points of view - a creator of a struct or interface and a user of it. It can be the same person. Let give me an example.

Imagine you have a task to send a HTTP request to an another service to store order data. You create an interface to abstract it. How should we call the interface? From the caller point of view, the external service is just a storage for your data.

```go
type OrderStorager interface {
    Store(o Order) error
}
```

From the instractructure point of view, it's not a storage. It's more like a regular HTTP call to somewhere else. The (pseudo) code can be:

```go
type AnotherMicroserviceService struct {
    client http.Client
}

func (ams AnotherMicroserviceService) Send(o Order) error {
    // build the HTTP request and send it
    // using ams.client.Do(req)
    return nil
}
```

When you work with the `OrderStorager` interface, is it important that it's a HTTP call? In my opionion - no.
On the other hand, we can call the interface `AnotherMicroserviceServices` and the method `Send()`. It will work but what if something changes and you want to store it in the service **and** in your local database at the same time. The interface would start lying about what it does. We subtly break open-closed principle from SOLID. Adding any cache or doing anything else what doesn't change the actual logic should be replacable and not visible from the user's point of view because it doesn't change the core behavior.

## Don't let your logic leak to the infrastructure

The infrastructure code should be as simple as possible. The reason behind it is simple: it's a hard to test area. Let me give you an example. Let's say we have a method in our repository.

```go
func (repo Repository) Get(ctx context.Context, id string) (record, error) {
    q := `SELECT id, name from TABLE where id = ? AND status = ?`
    row := repo.db.QueryRowContext(ctx, q, id, "active")
	r := record{}

	if err := row.Scan(&r.id, &r.name, &r.status); err != nil {
		return record{}, err
	}

	return r, nil
}
```

Can you find the problem? We subtly require that if the row has status different than `active` it's not allowed to fetch it. How can we test it? Only with real database running in the background. We can always mock the `sql.DB` but it's hacky for me. To fix it we can fetch the status from the database and check in in the service.

```go
var ErrRecordIsNotActive = errors.New("the record is not active")

func (s service) DoSomething(ctx context.Context, id string) error {
	r, err := s.repo.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("cannot do something: %w", err)
	}

	if r.status != "active" {
		return fmt.Errorf("cannot do something: %w", ErrRowIsNotActive)
	}

	// do the job
	return nil
}

func (repo Repository) Get(ctx context.Context, id string) (record, error) {
	q := `SELECT id, name, status from TABLE where id = ?`
	row := repo.db.QueryRowContext(ctx, q, id)
	r := record{}

	if err := row.Scan(&r.id, &r.name, &r.status); err != nil {
		return record{}, err
	}

	return r, nil
}
```

The refactoring made the testing of the logic very simple. All we need to do is to mock the repository and return an object with a proper status. No database running is required.

## Errors are part of our API

In the example above we has an another problem. When the row does not exists the `sql.NoRows` error leaks from the infrastructure to the service layer. We forget that errors are part of the API of our functions or packages. We should think about them as well to make the code more developer friendly. Let's think what should we do in the service to check if the record exist. We have two options: compare the error to `sql.NoRows` or compare strings. Both of them aren't very elegant solutions. To make it more handy, we refactor the repository once again to add a separate error.

```go
func (repo Repository) Get(ctx context.Context, id string) (record, error) {
	q := `SELECT id, name, status from TABLE where id = ?`
	row := repo.db.QueryRowContext(ctx, q, id)
	r := record{}

	if err := row.Scan(&r.id, &r.name, &r.status); err != nil {
		if errors.Is(err, sql.NoRows) {
			return record{}, ErrRecordDoesNotExist
		}
		return record{}, err
	}

	return r, nil
}
```

All we have to do is compare the error in the service to this one retrived from the repository.

## Keep interfaces close to the consumer of the API

I have a library that wraps the standard `sql.DB`. I prepared an interface with all functions from `sql.DB` struct to be 100% compatable with Go 1.15.

```go
type Database interface {
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Close() error
	Conn(ctx context.Context) (*sql.Conn, error)
	Driver() driver.Driver
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Ping() error
	PingContext(ctx context.Context) error
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	SetConnMaxIdleTime(d time.Duration)
	SetConnMaxLifetime(d time.Duration)
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	Stats() sql.DBStats
}
```

I implemented my requirements and gave other devs to use. After some time, I received a feedback that one team cannot use it. They received an error that their `sql.DB` doesn't have the `SetConnMaxIdleTime` function and the code doesn't compile. After a quick investigation I noticed that they use older Go version. The missing function was added to Go in v1.15. What's more, they reused the interface I provided so they couldn't change it. What did go wrong?

Their code had higher coupling with my library than it was needed. They reused the interface I provided. The fix for it was simple - copy the interface and remove the problematic function from the interface. When they was ready to upgrade the Go version they can just add the missing function back. Or not :)

This approach unlocks me to add more functions to my library without thinking if I'll break someone else's code. When the Go team adds a new function to `sql.DB` I can add it to the wrapper just after them.

But why does it work this way? The answer is simple - we invert the dependency. We don't require the library with the specific interface, we created our own that it's compatable with the wrapper. As long as the wrapper's interface is backword compatable, we can upgrade it without worries.

Keep the interface close to it's usage. Don't be afraid of coping. 

Didn't confinienced? Please take a look at [this interface](https://github.com/bkielbasa/gotodo/blob/ca0a2ba5035e665b2680ae6d62003cf041bfdc9d/repositories/postgres.go#L15).

```go
type PostgresRepository interface {
	GetToDo(id string) (string, bool, error)
	AddToDo(projectID, name string) (string, error)
	MarkToDoAsDone(id string) error
	MarkToDoAsUndone(id string) error
	ListToDos() (httpmodels.ListToDoResponse, error)
	ListProjects() (httpmodels.ListProjectsResponse, error)
	CreateProject(name string) (string, error)
	ArchiveProject(id string) error
}
```

The interface is quite big. Why did it happen? The developer (probably) definied the interface at the very beginning and just added more and more. When I want to use it and need only 1 single function from it I have to take everything together with it. I can, of course, create my own smaller interface to fix it. That's a good approach but why do I need the bigger interface then? It's redundant.

I explained more the topic in the article [is my interface too big](https://developer20.com/is-my-interface-too-big/) and recommend reading it.

## Don't overuse concurrency

The concurrency is awesome! I know but it's not trivial as well. We have a tendency to overuse cool tools and the concurrency is one of them. The very first thing to stress is that your code, in most cases, is already concurrent. If you're writing an HTTP service then every request is processed in a separate goroutine. The same with event handlers. That's not 100% true of course.

When do you need the concurrency, then? 

* when you have a proof it improves the performance. Yeah, that's right. You should write benchmarks to show you're right about your predictions.
* when you deal with multiple input/output sources like hard drive, network connection etc
* for batching some operations. I've used buffered channels to get send multiple redords to Kinesis in a single request. It decresed the network overhead a lot! Instead sending 100 messages I sent only one but bigger. In this case, the memory usage increased but you know... trade-offs :)

Remember that the code is easier when a goroutine has a single task to do. In some cases, limiting the concurrency can be a good ballance. Below, you can find an example that [ilustrates it](https://medium.com/@_orcaman/when-too-much-concurrency-slows-you-down-golang-9c144ca305a):

```go
package msort

import "testing"

var a []int

func init() {
    for i := 0; i < 1000000; i++ {
        a = append(a, i)
    }
}
func BenchmarkMergeSortMulti(b *testing.B) {
    for n := 0; n < b.N; n++ {
        MergeSortMulti(a)
    }
}

func BenchmarkMergeSort(b *testing.B) {
    for n := 0; n < b.N; n++ {
        MergeSort(a)
    }
}

/*
resutl:
BenchmarkMergeSortMulti-8              1    1711428093 ns/op
BenchmarkMergeSort-8                  10     131232885 ns/op
*/
```

In this benchmark, we run two goroutines for every step. It end up with **a lot** of goroutines that are fighting for resources. Please remember that goroutines are laightweight but are not for free. The Go scheduler has to schedule every single goroutine, manage them and so on. It takes time too. In this specific example, the concurrent version is 13 times slower than the regular one. In this article, Skyline AI limited to 100 goroutines by rewriting the alghoritm.


```go
var sem = make(chan struct{}, 100)

func MergeSortMulti(s []int) []int {
    if len(s) <= 1 {
        return s
    }

    n := len(s) / 2

    wg := sync.WaitGroup{}
    wg.Add(2)

    var l []int
    var r []int

    select {
    case sem <- struct{}{}:
        go func() {
            l = MergeSortMulti(s[:n])
            <-sem
            wg.Done()
        }()
    default:
        l = MergeSort(s[:n])
        wg.Done()
    }

    select {
    case sem <- struct{}{}:
        go func() {
            r = MergeSortMulti(s[n:])
            <-sem
            wg.Done()
        }()
    default:
        r = MergeSort(s[n:])
        wg.Done()
    }

    wg.Wait()
    return merge(l, r)
}
```

It helped a lot. The concurrent version is now 3 times faster than the original code. That's a huge improvement. I went one step further and changed the `sem` channel in every test and prepared tests for 1, 2, 5, 10, 20, 100, 1000, 10000, 10000 buffer size. Here are my results:

```
➜  mergesort go test -bench=. -benchmem
goos: darwin
goarch: amd64
BenchmarkMergeSortMulti1-8                    22          48911484 ns/op        164768355 B/op   1000009 allocs/op
BenchmarkMergeSortMulti2-8                    22          48692658 ns/op        164768528 B/op   1000012 allocs/op
BenchmarkMergeSortMulti5-8                    27          44530667 ns/op        164769609 B/op   1000037 allocs/op
BenchmarkMergeSortMulti10-8                   24          48077685 ns/op        164777560 B/op   1000205 allocs/op
BenchmarkMergeSortMulti20-8                   27          45076583 ns/op        164806822 B/op   1001415 allocs/op
BenchmarkMergeSortMulti100-8                  28          41923743 ns/op        165120146 B/op   1012978 allocs/op
BenchmarkMergeSortMulti1000-8                 19          58402295 ns/op        168198394 B/op   1125470 allocs/op
BenchmarkMergeSortMulti10000-8                 6         194675420 ns/op        181665217 B/op   1605549 allocs/op
BenchmarkMergeSortMulti100000-8                1        1015733662 ns/op        288276224 B/op   4085021 allocs/op
BenchmarkMergeSortMulti1000000-8               2         863970796 ns/op        250489984 B/op   4059603 allocs/op
PASS
ok      _/Users/bartlomiej.klimczak/Projects/mergesort  15.236s
```

We can see that the number of gourutines impacts the speed of the calculation as well as memory usage. On my computer, the most optimal value seems to be 100. The test is available below in this [gist](https://gist.github.com/bkielbasa/667b34ea77c0c4f1708d51b214cbdf03) if you want to play with it.


## Wrap errors with some context

Let's consider this simple pseudo-code in your repository.

```go
func (myRepo MyRepository) SomeOperation(arg type) error {
	if err := myRepo.doSomething(arg); err != nil {
		return err
	}

	return myRepo.doSomethingElse(arg)
}
```

The code looks simple and clear. That's cool. The problem can show up when we use this function more than once without adding any context to the error we received. What will happen when the `doSomething` function is used multiple times in the repository and just returns an error `cannot do something`? Will we know from which function call it comes from?

Think about wrapping errors as adding a stack trace to it. Make it meaningful and easy to reason about when you find it in your logs. The code above can be rewritten to a bit longer version using only the Go standard library.

```go
func (myRepo MyRepository) SomeOperation(arg type) error {
	if err := myRepo.doSomething(arg); err != nil {
		return fmt.Errorf("some operation: cannot do something: %w", err)
	}

	if err := myRepo.doSomethingElse(arg); err != nil {
		return fmt.Errorf("some operation: cannot do something else: %w", err)
	}

	return nil
}
```

We can make it easier thanks to `github.com/pkg/errors` package.

```go
func (myRepo MyRepository) SomeOperation(arg type) error {
	if err := myRepo.doSomething(arg); err != nil {
		return errors.Wrap(err, "some operation: cannot do something")
	}

	err := myRepo.doSomethingElse(arg)
	return errors.Wrap(err, "some operation: cannot do something else")
}
```

In the example above, the error message can be similar to this:

```
some operation: cannot do something: connection timeout
```

The only `connection timeout` says almost nothing what did go wrong. By adding the context to the error, we reproduce the path the request went and can follow it back.

## Simplify `if` statements

There's [Anti-if campaign](https://code.joejag.com/2016/anti-if-the-missing-patterns.html) and I like it. Today, I want to confince you to something different - removing `else` from your code. Too big conditions or too many of them can lead to complicated and spaghetti code. Nobody likes working with it, right? How can we avoid them? Let's consider the code below.

```go
func foo(bar string) (int, error) {
	var val int
	var err error

	if val, err := foobar(); err == nil {
		calc := 0
		for _, anotherVal := range barfoo(val) {
			err = anotherAction(val, anotherVal)

			if err == nil {
				calc++
			} else {
				return 0, err
			}
		}

		return calc, nil
	}

	retunr 0, err
}
```

The code can be refactored to:

```go
func foo(bar string) (int, error) {
	val, err := foobar()
	if err != nil {
		return 0, err
	}
	
	calc, err := calculate(barfoo(val), val)
	if err != nil {
		return 0, err
	}

	return calc, nil
}

func calculate(barfooval []int, val int) (int, error) {
	calc := 0
	for _, anotherVal := range barfooval) {
		err = anotherAction(val, anotherVal)
		if err != nil {
			return 0, err
		}

		calc++
	}

	return calc, nil
}
```

The second function looks more readable. The idea is simple - align the happy path to the left of the code and enclose complicated code with functions. 

{{< imgresize happy-path.png "500x500" "Alternate Text" >}}

If possible, get read of the `else` statement by returning early.

```go
if foo.OK() bool {
	err := bar()
	if err == nil {
		return true
	} else {
		return false
	}
}

return false

// after the change
if !foo.OK { // fliped
	return false
}

err := bar()
if err != nil { // fliped
	return false
}

return true

// or even simpler
if !foo.OK { // fliped
	return false
}

err := bar()
return err == nil
```

Remember one thing - optimize the code for reading. The readability improves the maintanance. Easy maintanance makes other devs happier.
