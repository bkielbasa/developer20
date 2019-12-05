---
title: "Golang Tips & Tricks #8 - Replacing mutexes with channel"
publishdate: 2019-12-09
categories: [Golang, Programming]
tags:
  - golang
  - modules
  - channel
  - mutex
---

Concurrency is simple in Go but it doesn't necessarily mean it's easy. Mutexes are one of the ways of synchronization. Go Mutex is a struct that allows multiple goroutines to share the same resources like memory or file access. This article doesn’t answer if channels should be used instead of mutexes. It depends on the [problem you’re solving](https://github.com/golang/go/wiki/MutexOrChannel). If the mutex is simpler and easier to read - go for it. Just keep in mind that they have [some issues](https://opensource.com/article/18/7/locks-versus-channels-concurrent-go) as well. Today, I want to show you how you can avoid using the mutex at all, in some scenarios.

If you want to limit access to the resource to only one goroutine, you can use channels instead of mutexes. Let's take a look at the code below.

{{< highlight go>}}
package main

import (
    "fmt"
    "sync"
    "time"
)

type handler struct {
    arr   []int
    mutex sync.Mutex
}

func (h *handler) handle(n int) {
    // some other logic

    h.mutex.Lock()
    defer h.mutex.Unlock()
    h.arr = append(h.arr, n)

    // some other logic
}

func main() {
    h := &handler{}

    for i := 0; i < 10; i++ {
        go func(n int) {
        h.handle(n)
        }(i)
    }

    time.Sleep(10 * time.Millisecond)
    fmt.Printf("values: %v", h.arr)
}
{{< / highlight >}}

Here we use mutexes to secure the shared memory between goroutines. There's another way of doing it - using a channel! The idea is to send the data to a channel and in another goroutine receive those data one by one and process. The behaviour is very similar to the previous one but without using mutexes at all. To do that we need a function where the processing will happen. We called it `startProcessing`.

{{< highlight go>}}
func (h *handler) startProcessing() {
    for n := range h.ch {
        h.arr = append(h.arr, n)
    }
}
{{< / highlight >}}

We have to change the `handle` function as well to send the data to the channel. Thanks to this change, we can process it in parallel.
{{< highlight go>}}
func (h *handler) handle(n int) {
    // some other logic

    h.ch <- n

    // some other logic
}

{{< / highlight >}}

The last thing we have to do is adding the channel to the struct and call `startProcessing` function after creating the instance of the handler.


{{< highlight go>}}
type handler struct {
    //...
    ch  chan int
}

//...

func main() {
    h := &handler{
        ch: make(chan int),
    }
    go h.startProcessing()
    //...
}
{{< / highlight >}}

And that's all! The whole code [is available in the playground](https://goplay.space/#_J8-NCgjUk_d). The program looks more clear and we didn't need any mutex at all. How about you? Do you prefer mutexes or channels? Let us know in the comments section below.

