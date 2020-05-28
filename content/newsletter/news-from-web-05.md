---
title: "News from the web #5"
publishdate: 2020-05-28
---

Hi there!

I hope you're well. This time, you'll learn about asynchronous preemption, inlining in Go and be able to read GW-BASIC source code.

## [OpenAI Model Generates Python Code](https://www.youtube.com/watch?v=fZSFNUT6iY8) #ai

Do you think that in the future computers will replace us? This is an example of code generation using OpenIA.

## [Microsoft Open-Sources GW-BASIC](https://devblogs.microsoft.com/commandline/microsoft-open-sources-gw-basic/)

Source code from 10th Feb 1983 is open sourced! If you know assembly or interested in archaeology - it's a perfect place for you. They give a brief background of the project and explain why they didn't use C etc for it.

## [Xbox and Windows NT 3.5 source code leaks online - The Verge](https://www.theverge.com/2020/5/21/21265995/xbox-source-code-leak-original-console-windows-3-5) #security

Microsoft published (without their will) another piece of code. The operating system is only used in a small number of systems worldwide so a source code leak isn’t a significant security issue but still :)

## [Houdini is a clone of Stockfish 8](https://groups.google.com/forum/#!topic/fishcooking/DygaIdBvJm0)

The guy had a simple idea: clone an open source project, change variable names to something "weird" and start selling it. The problem is that he violated the GPL.

## [Habits of High-Functioning Teams](https://deniseyu.io/2020/05/23/habits-of-high-performing-teams.html)

It's a list of good practices to keep your productivity as high as possible as the team. You'll learn why `git blame` isn't a good tool for you, why communication matters as well as how and why record only the present day.

## [How the biggest consumer apps got their first 1,000 users - Issue 25](https://www.lennyrachitsky.com/p/how-the-biggest-consumer-apps-got) #business

There are only a few well working tactics for getting the very first 1000 users. From this post, you'll learn how the biggest apps go their first costumers. It may help you in running your own business.

## [AssemblyScript: Passing Data to and From Your WebAssembly Program](https://www.jameslmilner.com/post/assemblyscript-passing-around-data/) #webassembly

WebAssembly is an interesting technology. There's still not that much materials about it. From this link, you'll learn how to pass basic data structures to WebAssembly like: numbers, arrays or strings.

## [Go: Asynchronous Preemption](https://medium.com/a-journey-with-go/go-asynchronous-preemption-b5194227371c) #golang

If you want to know a bit more about the scheduler in Go. In Go 1.14, the behaviour changed so it's worth reading materials like that.

## [Dgraph, GraphQL, Schemas, and CRUD](https://www.ardanlabs.com/blog/2020/05/dgraph-graphql-schemas-crud.html) #graphql #dgraph

Quite long article about Graphs and GraphQL. I'd say it's very basic and shows how to run a simple CRUD application using GQL and Dgraph but it can be a good starting point for experiments.

## [Three bugs in the Go MySQL Driver](https://github.blog/2020-05-20-three-bugs-in-the-go-mysql-driver/) #golang

Personally, I love articles like this one. It's very detailed but easy to understand summary of 3 bugs in Go MySQL driver that left there for a while.

## [Mid-stack inlining in Go](https://dave.cheney.net/2020/05/02/mid-stack-inlining-in-go) #golang

Dave gives us another cool post. This time about inlining in Go. Inlining is a compiler optimisation that allows reducing the overhead of calling functions by embedding small functions into their callers. If you want to know more about it - read it.

## [Diamond interface composition in Go 1.14](https://dave.cheney.net/2020/05/24/diamond-interface-composition-in-go-1-14) #golang

Go had a limitation on composing interfaces. Two interfaces couldn't be embedded into the same type if both of them contain the same function definition. Now it's fixed.

## [Good job Mojang, how did you even manage that?](https://www.reddit.com/r/softwaregore/comments/gqv7af/good_job_mojang_how_did_you_even_manage_that/) #security

Funny bug! Imagine you're uninstalling a game and... the whole hard disc is cleared! It reminds me another bug in [bumblebee](https://github.com/MrMEEE/bumblebee-Old-and-abbandoned/issues/123). [We are only humans](https://www.youtube.com/watch?v=r5yaoMjaAmE).

## [Go and CPU Caches](https://medium.com/@teivah/go-and-cpu-caches-af5d32cc5592) #golang #performance

Modern PCs have 3 level cache. Each of them has different access time and size. If you're working on high performance applications, knowing such stuff is a must. From this article, you'll learn how it works and how to take advantage of it in Go.

## [Robust gRPC communication on Google Cloud Run (but not only!)](https://threedots.tech/post/robust-grpc-google-cloud-run/) #golang #grpc

How to run gRPC app on GCP? Yeah, it can be easy. You'll learn how to do it with protobuf.

That's all I prepared for you this week. I hope you like it. Don't forget to let your friends know about the newsletter and see you next week!
