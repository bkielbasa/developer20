---
title: "Golang Tips & Tricks #6 - the _test package"
publishdate: 2019-08-14
categories: [Golang, Programming]
tags:
  - golang
  - tests

resources:
    - name: header
    - src: featured.png
---


Testing is one of the hardest stuff in programming. Today trick will help you organize your tests and the production code.

Let’s assume you have a package called orders. When you want to separate the package for tests from the production code you can create a new folder and write tests there. It will work but there’s a more clearer way - put your tests to the folder with you package but suffix the package’s name in tests with _test.

{{< highlight bash >}}
order.go # package orders
order_test.go #package orders_test
{{< / highlight >}}

I use this approach a lot and it helps to keep both prod and test code together but can can test the production code like from an external package. I hope you’ll find it useful.
