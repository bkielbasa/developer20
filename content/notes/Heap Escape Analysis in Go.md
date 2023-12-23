public:: true

- When the Go compiler cannot tell for sure that the variable under the interface won't be used longer than the caller function lives, it will move the variable to the heap instead of the stack.
- When we run our code with a parameter `-gcflags="-m"` we'll get (in the stdout) the heap escape analysis. The output will tell us what kind of compile optimisations are done during building a package.
- ```sh
  $ go tool compile -m demo.go
  demo.go:3:6: can inline demoFunction
  demo.go:8:6: can inline main
  demo.go:9:31: inlining call to demoFunction
  demo.go:4:9: moved to heap: data
  ```
- We can use it to finding out the reason why we have allocations on the heap rather on the stack.
- There can be many reasons why a variable escapes to the heap like:
   * [[Using interfaces in Go leads to moving to the heap]]