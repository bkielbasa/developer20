---
title: "Golang Tips & Tricks #1 - errors"
publishdate: 2019-03-04
categories: [Golang]
tags:
  - golang
---

You should use the package `github.com/pkg/errors` instead of `errors` package for errors in your applications. The default package lacks a few things like:

 * stack trace
 * easy appending message to the error
 * and more

It helps with debugging a lot. Below you can find an example error message with the stack trace.

![](/assets/posts/tipsandtrics01.png)


An important thing to remember is that you should wrap every error which is from any external library or your freshly created error you want to return.

{{< highlight go >}}
l, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", host, port))
if err != nil {
    return errors.Wrapf(err, "cannot start listening on port %d", port)
}
{{< / highlight >}}

or...

{{< highlight go >}}
return errors.Wrap(MyError{}, "could not handle the situation")
{{< / highlight >}}

If you want to add an extra message to already wrapped error, use WithMessage() function.
