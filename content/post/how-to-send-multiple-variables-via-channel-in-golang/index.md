---
title: "How to send multiple variables via channel in golang?"
publishdate: 2018-12-10
categories: [Golang, Programming]
tags:
  - golang
  - channels
  - concurrency

resources:
    - name: header
    - src: featured.png
---
Channels in golang are referenced type. It means that they are references to a place in the memory. The information can be used to achieve the goal.

Firstly, let’s consider using structs as the information carrier. This is the most intuitive choice for the purpose. Below you can find an example of a struct which will be used today.

```go
type FuncResult struct {
	Err error
	Result int
}

func NewFuncResult(result int) FuncResult {
	return FuncResult{Result: result}
}
```
The idea is to create a channel from the struct, pass the channel to a function and wait for the result.

```go
func funcWithError(r chan FuncResult) {
	r <- NewFuncResult(123)
}

func main() {
	r := make(chan FuncResult)
	go funcWithError(r)
	res := <- r
	if res.Err == nil {
		fmt.Printf("My result is %d!", res.Result)
	} else {
		fmt.Printf("The func returned an error: %s", res.Err)
	}
}
```

[Example on Go Playground](https://play.golang.org/p/t_ggprDWIXB)

Another solution is using functions in similar way to structs. This is more functional-programming way and may look less readable.

```go
func funcWithError(f chan func() (int, error)) {
	f <- (func() (int, error) { return 123, nil })
}

func main() {
	r := make(chan func() (int, error))
	go funcWithError(r)
	res, err := (<-r)()
	if err == nil {
		fmt.Printf("My result is %d again!", res)
	} else {
		fmt.Printf("The func returned an error: %s", err)
	}
}
```

[Example on Go Playground](https://play.golang.org/p/xXxYPuddJTw)

To simplify the code a bit, it is a good idea to define a custom type which will help keeping the code more readable.

```go
type FuncResult func() (int, error)

func funcWithError(f chan FuncResult) {
	f <- (func() (int, error) { return 123, nil })
}

func main() {
	r := make(chan FuncResult)
	//...
}
```
