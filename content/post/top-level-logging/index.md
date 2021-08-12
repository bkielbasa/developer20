---
title: "Top level logging"
publishdate: 2021-07-03
categories:
    - Golang
    - Programming
tags:
    - logging 
resources:
    - name: header
    - src: featured.jpg
---

I like having the core logic of our application free of distractions  like too many details or some "technical" details like logging or generating metrics. Of course, sometimes it's hard to avoid it. I found in many projects a situation where we put the logger very deeply inside of the code. At the end of the day, we had the logger almost everywhere. In tests, we had to provide the mocked implementation everywhere as well. In most cases, the logger is a redundant dependency. In this article, I'll argue that we should have the logger only in top level functions.

The idea behind the top level logging rule is simple - you log everything only in one place and don't pass the logger in lower layers of your application. What is the top level? For example, your CLI command, HTTP or event handler. Below, you can find an example with logging every error on handler level.

```go
type myHandler struct {
  logger log.Logger
  srv myService
}

func (h myHandler) operation(w ResponseWriter, r *Request) {
  body, err := io.ReadAll(r.Body)
  if err != nil {
    h.logger.Errorf("cannot read the body: %s", err)
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  req := request{}
  if err = json.Unmarshal(body, &req); err != nil {
    h.logger.Errorf("cannot read the body: %s", err)
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  err = h.srv.Operation(r.Context(), req.Param1, req.Param2)
  if err != nil {
    h.logger.Errorf("cannot execute the operation: %s", err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // return the success response
}
```

The code looks straightforward. I noticed we sometimes put the logger to into other places too. The `myService` can be a good example.

```go
type myService struct {
  logger log.Logger
}

func (s myService) Operation(ctx context.Context, param1, param2 int) error {
    result := myOperation(param1, param2)

    if result == 0 {
        // this shouldn't happen but when it does, we're ignoring such cases
        s.logger.Infof("the result is zero")
        return nil
    }

    // do some other operations

    if err := s.db.Persist(ctx, myCalculations); if err != nil {
        return fmt.Errorf("cannot persist X: %w", err)
    }

    return nil
}
```

We use the logger independently in the service level to let it become known about a potential corner case that's ignored. On the one hand, it makes sense. We don't want to return an error because our logic is prepared for such an edge case. On the other hand, we're doing two things:

* we add an uncesessary dependency to a service that doesn't really require it
* we're making this edge case harder to test

The last point may be the most controversial. How is it harder to test? All we have to do is provide values to param1 and param2 that will produce the result =0 and check if the method returns a nil. And yes, you'll be right. I showed you a simple case but imagine that if this statement is hidden somewhere deeper and to make sure that you're covering the right return nil case, you have to check it manually. What's more, someone may add another check **before** our target if statement. It may lead to a situation where our test still passes but it gives false information about which condition returns the nil.

```go
func (s myService) Operation(ctx context.Context, param1, param2 int) error {
    op := anotherCheck(ctx, param1)

    if op > threshold {
        return nil
    }

    result := myOperation(param1, param2)

    if result == 0 {
        // this shouldn't happen but when it does, we're ignoring such cases
        s.logger.Infof("the result is zero")
        return nil
    }

    // do some other operations

    if err := s.db.Persist(ctx, myCalculations); if err != nil {
        return fmt.Errorf("cannot persist X: %w", err)
    }

    return nil
}
```

In larger projects  you'll have more situations like this [depicted on source code]. Handling them this way hides some conscious decisions deeper in the code. What I can suggest in this example is creating a new error and return instead.

```go
var ErrEmptyResult = errors.New("the result is zero")

func (s myService) Operation(ctx context.Context, param1, param2 int) error {
    result := myOperation(param1, param2)

    if result == 0 {
        return ErrEmptyResult
    }

    // do some other operations

    if err := s.db.Persist(ctx, myCalculations); if err != nil {
        return fmt.Errorf("cannot persist X: %w", err)
    }

    return nil
}
```

Notice that we don't need the logger in the Operation method anymore. What about the log message? We can easily move it to the handler.

```go
func (h myHandler) operation(w ResponseWriter, r *Request) {
  // ...

  err = h.srv.Operation(r.Context(), req.Param1, req.Param2)
  if err != nil {
    if errors.Is(err, ErrEmptyResult) {
        // this shouldn't happen but when it does, we're ignoring such cases
        s.logger.Infof("the result is zero")
        return
    }

    h.logger.Errorf("cannot execute the operation: %s", err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // return the success response
}
```

The testing is going to be more understandable and precise. We're clearly saying what we're expecting from the method and be 100% sure about which `return` was called. The drawback is that the `if err != nil` statement in the handler may become very massive after the time. That can happen, of course but in such cases, I'd considered if the handler or the logic in this place is too big and it may be worth splitting it into smaller parts.

[^1]: [Aspect-Oriented programming](https://en.wikipedia.org/wiki/Aspect-oriented_programming) is a good answer for it but in Go it may be challenging to introduce it.
