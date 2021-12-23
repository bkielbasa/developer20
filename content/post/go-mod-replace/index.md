---
title: "`replace` directive in go modules"
publishdate: 2021-12-15
categories: 
    - Golang
    - Programming
tags:
  - golang
  - modules

resources:
    - name: header
    - src: featured.png
---

Sometimes, we may want to use a library but a slightly modified version. It happens very often when we develop the library but test it in the context of an application. Go has a handy mechanism in go modules that can help us with it.

To make it work, we have to clone the library somewhere near the target project and run the following command in the application's folder.
```sh
go mod edit -replace github.com/my/library ../path
```

The path can be both relative (to the application root folder) or absolute. The `go.mod` file will be edited as follows.

```
module github.com/myorg/app

require (
	github.com/thirdpart/library v1.1.0
)

replace github.com/myorg/library => ../path
```

Notice that you can modify the `go.mod` file without using the the `go mod edit -replace` command.

From this moment, every time we compile the application, the updated dependency will be used. There's only one thing to remember. When you finish working, don't forget to remove the `replace` directive from your `go.mod` file.

Another usage of the `replace` directive is when we want to replace one library with its (maybe our) fork. Unfortunately, it happens that a library author stops maintaining the code but we need to make some changes to it. Replacing the dependency may be the easiest way of solving this problem.

```
module github.com/myorg/app

require (
	github.com/thirdpart/library v1.1.0
)

replace 	github.com/thirdpart/library => github.com/myorg/thirdpart-forked
```

The mechanizm is simple and very powerful. Hope you'll find it helpful!