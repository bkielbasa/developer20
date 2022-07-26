---
title: "Honestly about why Go sucks (or not)"
resources:
    - name: header
    - src: featured.png

toc: true
---

Go is very opinionated. There are arguments that are based on personal preferences like "I don't like the syntax" and much more specific. In this article, I'll focus on the second type of arguments why Go isn't the best language and confirm/denied them. My goal is to tell you the truth about the language.

{{< table_of_contents >}}

## Arguments agains the language

### Lack of Function Overloading and Default Values for Arguments (https://www.toptal.com/go/4-go-language-criticisms)

Yes, Go doesn't have those features. And probably will never have. The argument here is that developers have to write more code than they have to. Right now, we have to write functions like this:

```go
func (wd *remoteWD) WaitWithTimeoutAndInterval(condition Condition, timeout, interval time.Duration) error {
    // the actual implementation was here
}

func (wd *remoteWD) WaitWithTimeout(condition Condition, timeout time.Duration) error {
    return wd.WaitWithTimeoutAndInterval(condition, timeout, DefaultWaitInterval)
}

func (wd *remoteWD) Wait(condition Condition) error {
    return wd.WaitWithTimeoutAndInterval(condition, DefaultWaitTimeout, DefaultWaitInterval)
}
```

Instead of just one single method with default values

```go
function (wd *remoteWD) Wait(condition, timeout = DefaultWaitTimeout, interval = DefaultWaitInterval) {
    // actual implementation here
}
```

I partially agree with this one. Sometimes it's useful to have some default values but it has disadvantages too. Let's consider the situation where we want to call the method `Wait()` but only pass the `interval` parameter. Default values won't help here. We'll have to write another method that accepts only the `interval` parameter or call `Wait()` with the default value provided explicitly.

```go
wd.Wait(DefaultWaitTimeout, myInterval)
```

That's not a good developer experience. The problem will become bigger if we have more parameters with default values. But, we can solve the problem with already existing syntax! The first way of doing it is using variadic functions. At the end of a function, we can provide a list of parameters of a specific type. Inside of the body, we can access them as a regular slice.

```go
func sum(nums ...int) {
    fmt.Print(nums, " ")
    total := 0
    for _, num := range nums {
        total += num
    }
    fmt.Println(total)
}
```

The caller of the function can provide as many arguments as he wants.

```go
sum(1,2,3,4,5,6,7,8,9,10)
```

There's one requirement - the type must be the same. Trailing argument has to be last in function arguments but not necessary the only one.

```go
func sum(w io.Writer, nums ...int) {
	fmt.Fprint(w, nums, " ")
	total := 0
	for _, num := range nums {
		total += num
	}
	fmt.Fprintln(w, total)
}
```

We can use this feature in our example by passing a list of functions that edit the internal state of a struct as shown below.

```go

func main() {
	s := newServer(withPort(8081))
	fmt.Print(s)
}

type server struct {
	port    int
	timeout time.Duration
}

func newServer(opts ...option) *server {
	s := &server{
		port:    8080, //default port
		timeout: time.Second,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

type option func(s *server)

func withPort(p int) func(s *server) {
	return func(s *server) {
		s.port = p
	}
}

func withTimeout(tm time.Duration) func(s *server) {
	return func(s *server) {
		s.timeout = tm
	}
}
```

Of course, we have to write more code but, on the other hand, please notice how flexible our code is becoming. We can add `error` in the `option` function and add validation to it.

```go
type option func(s *server) error


func withPort(p int) func(s *server) error{
	return func(s *server) error {
        if port < 80 {
            return fmt.Errorf("cannot provide port number lower than 80, given %d", p)
        }
		s.port = p
        return nil
	}
}

```

There's a similar pattern that's shown below.

```go
wd = wd.WithInterval(myInterval)
wd = wd.WithCondition(condition)
```

Depending on your use case, you can make the type you're working on mutable or immutable. What's more, this pattern leads to simpler and smaller functions what's always a benefit.

To sum up, in Go we don't have function overloading or default arguments but I don't feel it's needed. We can achieve a similar goal with existing syntax. We'll end up with more lines of code but with a more flexible one. Maybe you won't like this answer but I disagree that it's a huge Go's disadvantage.

**Declined**

### Lack of Generics (https://www.toptal.com/go/4-go-language-criticisms)

We have generics since [1.18](https://go.dev/blog/go1.18). Until then, the argument was valid. I didn't feel the need for generics myself but it looks like there's a huge number of people that needed that.

To be fair, before the 1.18 release we did have generics in Go. One example is the method `make()`. The problem is that we couldn't create our own generic functions or types. Because of that, to make our code generic we had to use `interface{}` and validate the type in runtime. As you may guess, it can lead to some mistakes that lead to a runtime panic. A good example of such a function is [map/reduce and filter](https://www.reddit.com/r/golang/comments/1m25a1/map_reduce_and_filter_in_go/) operations.

It's not the case anymore. Of course, we have to wait sometime to have [libraries](https://pkg.go.dev/golang.org/x/exp/slices) in the standard library but it's now just a matter of time. 

**Declined**

### Dependency Management (https://www.toptal.com/go/4-go-language-criticisms)

Yeah, it was a big issue in the Go ecosystem but it's not a problem at all since [go 1.11](https://pkg.go.dev/cmd/go#hdr-Modules__module_versions__and_more). I could stop here but there's something to add.

Firstly, we had to wait years for proper dependency management. Till then, all projects shared the same folder where all dependencies were downloaded. That led to problems like working with multiple projects with different vendors' versions.

After some time, the Go Team created an experiment library called [dep](https://github.com/golang/dep). In the meantime, the community wrote their tool that tried to fix the same problem. When the Go Team introduced go modules we had:

* go modules
* dep
* godep
* glide

... and many more. Some people migrated from third-parted dependency management to `dep` because they thought that is the "official" one. After some time, they had to migrate once again to go modules.

I've noticed that the situation created a lot of confusion and sometimes frustrations in the community. Including me.

**Declined**

### Not very expressive (https://raevskymichail.medium.com/why-golang-bad-for-smart-programmers-4535fce4210c#80de)

The point here is, IMO, that in Go we produce more code compared to Dlang. The author provides the same application written in Go and Dlang. Firstly, the Go code can be simplified a bit by removing the usage of the packages `flag` and `bufio`.

```go
package main

import (
    "fmt"
    "log"
    "os"
    "io"
)

func main() {
    var r io.Reader

    if len(os.Args) > 1 {
        var err error
        r, err := os.Open(os.Args[1])
        if err != nil {
            log.Fatal(err)
        }
        defer r.Close()
    } else {
        r = os.Stdin
    }

    text, err := io.ReadAll(r)
    if err != nil {
      log.Fatal(err)
    }

    fmt.Print(string(text))
}
```

Both `os.File` and `os.Stdin` implements `io.Reader` interface so we can use this to make the code shorter.
We reduced the number of lines from 44 to 30. That's not that bad. The biggest difference between Go's code and the one in `D` is how the `if` statement looks like and how both languages handle errors. Is Go's version less readable than the one in `D`? I'm biased so you tell me.

**No judgment**

### Lack of stack traces in the errors

> Stacktraces are possible, but they have to be handrolled in the error handling. You need to make a log statement before returning an error, and in that log you need to make sure the line number, function name, file name, etc... are all used.

The argument is that when you create an error in Go you don't get any information about where the error was returned or created. I found myself debugging where the specific error comes from and sometimes it's very difficult to do so. To achieve something similar, you have to wrap errors yourself like this:

```go
err := sendOrder(arguments)
if err != nil {
  return fmt.Errorf("cannot complete the operation: %w", err)
}
```

You can overcome this inconvenience by wrapping the errors into some telling and unique comments. On the other hand, not everyone, and not always, has the self-discipline to do it conscientiously.

What's more, people who are just learning this language simply do not know it and (which often happens) even ignore errors. Not to mention wrapping them properly. Sometimes, it's painful to me even now. There's a package that possibly fixed that but it's [https://github.com/pkg/errors](archived) and it doesn't seem to be a preferred way of doing it.

**Confirmed**

### It's standard library isn't "All you need"

I totally agree with it. I've complained about it [many times](https://twitter.com/kabanek/status/1513794929931472900). The standard library is missing things like:

* working with YAML files
* the standard `log` package is extremely limited. It's perfect for very simple apps but when you want to add logging in more complex app, you'll need something like [logrus](https://github.com/sirupsen/logrus) or [zerolog](https://github.com/rs/zerolog).
* people who write bigger CLI tools often use [cobra](https://github.com/spf13/cobra) or [urfave/cli](https://github.com/urfave/cli) because the standard one `flag` package is too limited in some areas.
* the standard [mux](https://pkg.go.dev/net/http#ServeMux) doesn't support regexps/pattern matching and named routes
* and many other

There are good parts of the stdlib but it's not definitely "all you need[^allyouneed]" when building a standard application/cli tool.

**Confirmed**

### The lie that it is more performant than Java or C#

I love this argument because depending on what you want to prove, you can have a different result. I've done myself two benchmarks when comparing Go to Java. In one of those tests, Java was about 10% faster because the JIT did so great work. Of course, the cold start was bigger but after some requests, the Java app was faster than the same written in Go.

In one of the benchmarks. In another one, I had the opposite results. In both experiments, I've tried to test different parts of languages.

People say that Go is fast. It is but in some areas. In others, Go can be [one of the slowest (slower than python or PHP)](https://github.com/mariomka/regex-benchmark). Can we say objectively say that one language is faster than another one? I don't know if there's a single benchmark to answer the question.

**No judment**

### `nil` and type safety

Another argument against Go is the fact that Go has `nil` and has weak type safety. Let's consider the following code:

```go
func createInvoice(params createInvoiceParams) (*invoice, error) {
    // the actual impl
}

invo, err := createInvoice(params)
if err != nil {
    return fmt.Errorf("cannot complate the operation: %w")
}

fmt.Print(invo.ID())
```

The `createInvoice` method returns a pointer to the invoice and an error. If there's no error, we very often assume that the `invoice` won't be a `nil`. **We assume**. The only way of making sure the `invo` variable isn't a `nil` is by explicitly checking it. That's a bit problematic. If we won't do it we may see, at some point, well-known panic `runtime error: invalid memory address or nil pointer dereference`.

We can make sure that the `invoice` isn't a `nil` if we change the return type from the pointer to the value type.

```go
func createInvoice(params createInvoiceParams) (invoice, error) {
    // the actual impl
}
```

It solves the problem of panic but introduces another one. What if the function returns a `nil` error but a default (not initialized) instance of the `invoice`? We may not check it or notice this fact. This situation may lead to even more hard-to-find bugs.

Do you think this problem doesn't happen in real-world apps? I produced a bug like that at least a few times. What's more, the `http.Client.Do` method may return a non-nil response with a non-nil error. It may lead to a gotcha described in [50 Shades of Go](http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/index.html#close_http_resp_body).

If your function accepts an interface you cannot be 100% sure if there isn't a `nil`. Adding the `nil` checks in every method sounds crazy. You can find many places where we throw a panic in such cases in the [stdlib](https://cs.opensource.google/go/go/+/refs/tags/go1.18.4:src/context/context.go;l=435).

Rust, on the other hand, checks situations like that at compile time. It means if you're not making some crazy things you should be free from mistakes like the one above.

**Confirmed**

### Go isn't OOP

Go doesn't have classes, abstract methods or inheritance. It doesn't mean you cannot use object-oriented programming in this language. Before arguing about it we have to understand [what OOP is](https://medium.com/@egonelbre/relearning-oop-89f10e0e2f68) and remember that it's a [paradigm](https://medium.com/@egonelbre/paradigm-is-not-the-implementation-af4c1489c073). 

Here are some articles for further reading:

* https://flaviocopes.com/golang-is-go-object-oriented/ - the author answers the question if Go is Object-Oriented programming language
* https://www.toptal.com/go/golang-oop-tutorial - na example how you can use OOP in Go in practice.

The answer to the question can be only one:

**Declined**

### Error handling

Some people complain about it. Some people love it. The truth is somewhere in between. Let me explain.

In Go, the [recommended way of handling any errors](https://go.dev/doc/effective_go#errors) is using the `errors` package from the standard library. To create a new error we can use the `errors.New()` method or `fmt.Errorf()`. We can compare errors using `errors.Is()` and `errors.As()` functions. Pretty straightforward.

However, if we want to add a stack trace to the error message to see where the problem occurred we have to use the third-party library. Be designed, the error message should be our stack trace. It means, we can write one error with another and add additional information to it.

```go
// inside of `confirmAccount` func
err := activateUser(ctx, userID)
if err != nil {
  return fmt.Errorf("cannot confirm the account: %w", err)
}
```

To wrap the error with another I used the `%w` directive that stores the original error within the new one. It doesn't sound complicated, does it?

When I log the error I'll see something like that.

```sh
cannot confirm the account: the account has been already activated
```

Looks clear and elegant, doesn't it? Much more helpful than a very verbose stack trace.

```sh
goroutine 1 [running]:
main.Example(0x2080c3f50, 0x2, 0x4, 0x425c0, 0x5, 0xa)
        /Users/bill/Spaces/Go/Projects/src/github.com/goinaction/code/
        temp/main.go:9 +0x64
main.main()
        /Users/bill/Spaces/Go/Projects/src/github.com/goinaction/code/
        temp/main.go:5 +0x85

goroutine 2 [runnable]:
runtime.forcegchelper()
        /Users/bill/go/src/runtime/proc.go:90
runtime.goexit()
        /Users/bill/go/src/runtime/asm_amd64.s:2232 +0x1

goroutine 3 [runnable]:
runtime.bgsweep()
        /Users/bill/go/src/runtime/mgc0.go:82
runtime.goexit()
        /Users/bill/go/src/runtime/asm_amd64.s:2232 +0x1
```

Oh the other hand, it's very difficult to keep writing good comments for errors. You have to add it in every place that makes sense. Maybe probably everywhere? It's challenging to keep the discipline to do it right. It's getting harder when you want to be consistent across the team.

**Confirmed**

## Summary
Go is far from being perfect. It has some pros and cons. I use it in everyday job and it works for me. There are some areas where Go won't be a good fit in some type of projects. That's why we have Rust, TypeScript, Python, Lua and more.

There are so many choices so if you don't like Go's philosophy, the syntax or anything else - there's so many other options you can choose. Please remember that the grass is always greener on the other side of the fence.

[^allyouneed]: Do you know a programming language that contains all you need for everyday work?