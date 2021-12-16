---
title: "Top level logging"
publishdate: 2021-12-16
categories:
    - Golang
    - Programming
tags:
    - logging 
    - golang
resources:
    - name: header
    - src: featured.jpg
---

I like having the core logic of our application free of distractions like too many technical "details" like logging or generating metrics. Of course, sometimes it's hard to avoid it. I found in many projects a situation where we put the logger very deeply inside of the code. At the end of the day, we had the logger almost everywhere. In tests, we had to provide the mocked implementation everywhere as well. In most cases, the logger is a redundant dependency. In this article, I'll argue that we should have the logger only in top-level functions.

The idea behind the top-level logging rule is simple - you log everything only in one place and don't pass the logger in the lower layers of your application. What is the top-level? For example, your CLI command or an HTTP or event handler. Below, you can find an example of logging every error on the handler level.

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

The code looks straightforward. I noticed we sometimes put the logger to into other places, too. The `myService` can be a good example.

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

We use the logger independently at the service level to let it become known about a potential corner case that's ignored. On the one hand, it makes sense. We don't want to return an error because our logic is prepared for such an edge case. On the other hand, we're doing two things:

* We add an unnecessary dependency to a service that doesn't require it
* We make this edge case harder to test

The last point may be the most controversial. How is it harder to test? All we have to do is provide values to `param1` and `param2` that will produce the `result = 0` and check if the method returns a `nil`. And yes, you'll be right. How can you make sure that the test passes because `result = 0`? You can do it in a few ways:

1. check the code coverage - if those lines are green, we've done it. The problem will happen when someone will update the code **before** our target if statement. It may lead to a situation where our test still passes but it gives false information about which condition returns the `nil`.

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

2. Use the debugger to make sure those lines are executed - the drawback is similar to the previous idea
3. Mock the logger and check logged messages. That will work quite well. No way for misunderstanding in tests and so on. For me, it's a bit hacky, don't you think?

You’ll have more situations like this (depicted on source code). Handling them this way hides some conscious decisions deeper in the code. What I can suggest in this example is to create a new error and return it instead.

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

Notice that we don't need the logger in the `Operation()` method anymore. What about the log message? We can easily move it to the handler.

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

The testing is going to be more understandable and precise. We're clearly saying what we're expecting from the method and be 100% sure about which `return` was called. The drawback is that the `if err != nil` statement in the handler may become very massive after a time. That can happen, of course. In such cases, I'd consider if the handler or the logic in this place would be too big and it may be worth splitting it into smaller parts.

## No more logs in other places?

What I'm trying to do is to convince you to avoid using logger in deeper layers of your code. There may be situations that it may be hard. On the other hand, having the logger may be useful. One of usage that comes to my mind is letting know about some edge cases as showed above but hidden deeper in the code. Another one is adding logs with trace or debug level and enable the proper log level when we start experience weird problems on production.

Your usage may be valid, of course. The problem is when we overuse the logger and use it when we have too complicated code or our tests cover too much code at the same time and it's hard to find out where's the root cause. Logging shouldn't be a replacement for refactoring. I mean, it can be and can be beneficial in short term. In longer term, it may only cover technical depts by introduciong another one.