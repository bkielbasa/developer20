---
title: "Golang Tips & Tricks #4 - internal folders"
publishdate: 2019-03-25

---
While developing a library, we create a directory structure to keep the code organized. However, some exported functions or struct should not be used by users of the library. The achieve that, call the package `internal`.

```
.
├── bar
│   ├── bar.go
│   └── internal
│       └── foobar.go
├── internal
│   └── foo.go
└── main.go
```

In the example above, the `foo.go` can be included only in the `main.go`. What's more, only the `bar.go` is able to include the `foobar.go` file. It means, that only the direct parent of the `internal` package is allowed to use it's internals.

Use this feature reasonably!
