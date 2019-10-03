---
title: "GoCracow #5 Golang and BDD - Bartłomiej Klimczak"
publishDate: 2019-09-10
---

<iframe width="720" height="405" src="https://www.youtube.com/embed/G-wWOitDZAU" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

There is godog library for BDD tests in Go. I found this library useful but it run as an external application which compiles our code. It has a several disadvantages:

* no debugging (breakpoints) in the test. Sometimes it’s useful to go through the whole execution step by step
* metrics don’t count the test run this way
* some style checkers recognise tests as dead code
* it’s impossible to use built-in features like build constraints.
* no context in steps - so the state have to be stored somewhere else - in my opinion, it makes the maintenance harder

In this presentation I tell about my library which should solve those issues.