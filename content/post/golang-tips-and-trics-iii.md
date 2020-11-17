---
title: "Golang Tips & Tricks #3 - graceful shutdown"
publishdate: 2019-03-18
categories: [Golang, Programming]
tags:
  - golang
  - graceful shutdown

---
In the microservices' world, one thing what's worth considering is a graceful shutdown. This is important to not lose data while shutting down a container. The container orchestrator like Kubernetes can restart the container by sending `SIGTERM` or `SIGINT` signal. Those signals can be handled to safely close all connections and finish background tasks.

Signals are propagated using `os.Signal` channel. You can add the above code to your main.

```go
var gracefulStop = make(chan os.Signal)
signal.Notify(gracefulStop, syscall.SIGTERM)
signal.Notify(gracefulStop, syscall.SIGINT)
```

Then, we need a goroutine to handle signals.

```go
go func() {
       sig := <-gracefulStop
       // handle it
       os.Exit(0)
}()
```

If we serve the HTTP server, the first thing we can do is [shutdowning](https://golang.org/pkg/net/http/#Server.Shutdown) the server.

```go
go func() {
       sig := <-gracefulStop
       server.Shutdown(ctx)
       // handle it
       os.Exit(0)
}()
```

It will shut down the server without interrupting any active connections. But what about the background tasks? There are, at least, 3 approaches I found which solves the problem.

## Wait!

You can add a `time.Sleep(2*time.Second)` statement and just exit. I personally don't like the solution because some tasks may need more than `X` seconds. On the other hand, setting to high sleep time is not a good idea eather. This is definitely the easiest way of doing it.

## Use channels

An another way you can achieve the goal is using channels. Here's how it works: you create two channels. The first one will communicate tell goroutines that it's time to stop and the second that time's up and we're exiting.

```go
var closing = make(chan struct{})
var done = make(chan struct{})

// pass both channels to background processes

go func() {
       sig := <-gracefulStop
       closing <- struct{}
       time.Sleep(2*time.Second)
       done <- struct{}
       os.Exit(0)
}()
```

Thank's to this, all background tasks have 2 seconds to finish up their work and then we exit.

## Wait groups

The `sync` package has the `WaitGroup`. The `WaitGroup` waits for a collection of goroutines to finish. The idea behind it is to create the `WaitGroup` in the main file and pass it to background tasks. And after that, we call `Wait()` function after receiving the closing signal.

```go
wg := sync.WaitGroup{}
var closing = make(chan struct{})

// pass both wait group to background channels

go func() {
       sig := <-gracefulStop
       closing := <- struct{}{}
       wg.Wait()
       os.Exit(0)
}()
```

Using this technique, we send the signal on the `closing` channel that those processes should stop their work and using `wg.Wait()` to wait when it will happen. When those background processes won't exit, the container orchestrator will terminate them anyway.

If you know other approaches, just leave a message in the comment's section.
