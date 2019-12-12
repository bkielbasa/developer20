---
title: "I want to learn Go - how to start?"
publishdate: 2019-12-09
resources:
    - name: header
    - src: featured.jpg
categories:
    - Golang
tags:
  - golang
---

You can find a lot of materials about Go (including this blog) but it's hard to find the best place to start. This article's goal is to sum up the most valuable materials I found to help others. I focus only on free materials.

## Fundamentals

* [The Go tour](https://tour.golang.org/welcome/1) - IMO the absolutely must do. You can try Go without installing it. You'll learn some basic syntax and concepts step by step,
* [Go by example](https://gobyexample.com/) - if you're confused how to use a certain part of the language it's possible you'll find an example of it on this page. Extremely usefull.
* https://golangnews.org/ and http://www.go-gazette.com/ newsletters - it's a great source of good quality materials related to Go,
* [50 Shades of Go: Traps, Gotchas, and Common Mistakes for New Golang Devs](https://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/) - a list of "gotchas" in Go. It's a long list and there's now way you'll remember everything but it's a good reference for future.
* [How to Write Go Code](https://golang.org/doc/code.html) - some people recommend it but actually I don't. It uses `GOPATH` and doesn't use [go modules](https://blog.golang.org/using-go-modules) what makes it a bit outdated.
* [Golang Tutorial for Beginners](https://www.youtube.com/watch?v=YS4e4q9oBaU) [video] - very basic introduction to the language but with quite good explenation of some context you may not catch while finishing The Go tour. Recommended as follow-up.
* [Learn Go in 12 Minutes](https://www.youtube.com/watch?v=C8LgvuEBraI) [video] - covers the very basic Go aspects. You definitely won't know how to write Go code after watching it but it's a good demo of the syntax and a few Go concepts.
* [gocode.io](https://www.gocode.io/) - if you like gamification then it can be a fun way of starting playing with the language.

## After that...

* [Effective Go](https://golang.org/doc/effective_go.html) - it's a list of tips for writing clear, idiomatic Go code. It's the gold standard and definitely must read.
* [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) - as authors say - it's a suplement to Effective Go. It's an official list of common mistakes with explanation.
* [Go blog](https://blog.golang.org/) - if you want to be up-to-date with changes to the language and its tools
* [Ardan labs blog](https://www.ardanlabs.com/blog/) - one of the best Go blogs. You'll find a series of post about the GC or go modules and much more.

## Books

* [The Little Go Book](https://www.openmymind.net/The-Little-Go-Book/) - a bit old book which covers fundamental aspects of Go. It doesn't contain anything special. I'd rather think about it like the supplement of The Go tour.
* [An Introduction to Programming in Go](http://www.golang-book.com/books/intro) - much more complex book which tells you more about concurrency or testing (veeery basic). It contains simple "problems" to solve which will help you understand those aspects better.
* [Go Bootcamp](http://www.golangbootcamp.com/) - basicaly it covers the same topics as the book above but in more details. It contains links to external resources like [Rob's talk](https://vimeo.com/49718712) what's a huge plus.
* [Webapps in Go](https://leanpub.com/antitextbookGo) - If you want to learn how to write a web app - the book is for you. You'll find there how to connect to a DB, how to work with forms, upload a file etc.
* [Test-driven development with Go](https://leanpub.com/golang-tdd) - writing tests in prev materials isn't covered very well. In some cases, only in a few sentences. Unfortunatelly, it's not finished yet.
* [Build web application with golang](https://github.com/astaxie/build-web-application-with-golang) - a book translated to many languages where you'll learn how to build a web application in Go using [beego](https://beego.me/).

As you may guess, the best teacher is practice! I hope that the list will help you finding useful materials and will speed your learning up! If you know any other cool materials, let me know in the comments section below. Cheers!
