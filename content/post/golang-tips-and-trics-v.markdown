---
title: "Golang Tips & Tricks #5 - blank identifier in structs"
publishdate: 2019-07-22
categories: [Golang, Programming]
tags:
  - golang
  - strucs

---
While working with structures, there's a possibility to initialize the structure without providing the keys of fields.

```go
type SomeSturct struct {
  FirstField string
  SecondField bool
}

// ...

myStruct := SomeSturct{"", false}
```

If we want to force other (or even ourselfs) to explicitly providing the keys, we can add `_ struct{}` in the end of the structure.


```go
type SomeSturct struct {
  FirstField string
  SecondField bool
  _ struct{}
}

// COMPILATION ERROR
myStruct := SomeSturct{"", false}
```

The code above will produce `too few values in SomeSturct literal` error. [Try it yourself](https://goplay.space/#aq8-_U65YKx).

In Go, the blank identifier in the type has size `0` and it means "blank identifier." Because it's black, it cannot be accessed directly.
It forces us to use the fields' names explicitly. This technique can be used to avoid bugs that might arise from specifying the arguments in an incorrect order.
