---
title: GoGoConf 2019 - report
publishdate: 2019-07-05
tags:
  - golang
  - conference

resources:
    - name: header
    - src: featured.png
---
Recently I've been on the GoGoConf conference in Cracow. It was a cool opportunity to learn more and meet interesting persons. Today I'll tell you about my thoughts regarding every talk from 2019 edition. Most of the talks don't have video available yet but when the videos will be published I'll update the post.

## Tackling contention: the monsters inside the ’sync.Locker’ - Roberto Clapis

I personally like Roberto a lot for the way he behaves and how professional he is. In his talk he explained how to fix performance problems using tools like pprof and trace. One of my takeaways from this talk is that we should use atomic package only really when necessary. It's not worth use the atomic package instead of mutex because the mutexes use the atomic anyway.

[![Tackling contention: the monsters inside the ’sync.Locker’](https://img.youtube.com/vi/ok4NEfqAXb0/0.jpg)](https://www.youtube.com/watch?v=ok4NEfqAXb0)

## Patterns for effective and effortless Observability - Prakash Mallela

Prakash raised an important topic - how to improve the observability in our project without or with little effort. This talked made me think a lot about how I observe my microservices and how I can improve this. It wasn't said loud we should automate it as much as we can. We can achieve that buy providing a custom mini-framework with configured middlewares which will ,for example, automatically log all slow queries and save somewhere execution times to external dependencies (database, event bus, external storage etc).

[![Tackling contention: the monsters inside the ’sync.Locker’](https://img.youtube.com/vi/fwrqVwuQE10/0.jpg)](https://www.youtube.com/watch?v=fwrqVwuQE10)

## ALEPH: scaling and speeding blockchain up with DAG’s and Go - Michał Świętek

What's cool about Michał's talk is that they showed how they do research about algorithms they use. Michał explained why the switched from python to Go and how efficient it was for them. I know just basics about blockchain so this part wasn't that interesting to me. The talk was full of specific data and facts.

## Packages & modules - Oleg Kovalov

After the presentation I had a small talk with a guy who uses Go in his pet project. He said that he has problem with organizing code in the project so the topic is needed. Oleg explained how he organizes the code in his projects and explained why. To be honest, I don't fully agree with him, mostly regarding [Clean Architecture](https://www.amazon.com/Clean-Architecture-Craftsmans-Software-Structure/dp/0134494164). I think I'll write an article about how I do that and why. If you want me to write the article - let me know in comments section below.

## Developing a Go API client: the do’s and don’ts - Anthony Seure

Have you ever written a library used by someone else except you? If yes, the talk is for you! :D Writing libraries is hard because you have to face many challenges like changes in the API. How to make your API extendable and usable for others? Anthony gives a few tips how it can be achieved, talks about pros and cons of some solutions.

## Diagnose your Golang App anytime anywhere! - Mateusz Dymiński

This talk is very similar to Roboerto's but he goes more deeper and showes in real-time demos how to use pprof to analise applications. In this talk, we can learn how to attach a debugger to the running application. On the one hand, I don't think we showed much more that we can read from official [blog post](https://blog.golang.org/profiling-go-programs), but on the other hand it was awesome to see it in practice.

## Final words

As you can see, presentations on the conferences touched various of topics. Some of them loosely related with Go. If you have any other awesome talks from other conferences, just send the link in the comments section below. Which other Go conferences are worth attending?
