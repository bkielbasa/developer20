---
title: "Enums in Go"
publishdate: 2022-10-15

categories:
  - Golang

resources:
  - name: header
  - src: featured.png
---

Go doesn't support enums. You probably know that. However, we can, somehow, simulate them using a few techniques or tools. Today, I'll tell you how I use enum-like types and what are other possible options to improve it.

When you start reading about enums in Go you'll quickly find a tip to use `iota`.

```go
const (  // iota is reset to 0
        c0 = iota  // c0 == 0
        c1 = iota  // c1 == 1
        c2 = iota  // c2 == 2
)
```

In every `const` block its value is reset to `0` and increment when declaring another constant as shown above.