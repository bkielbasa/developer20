---
title: "News from the web #1"
publishdate: 2020-04-29
---

Hi there!

This newsletter is a small experiment. I want to check if there are people who like such summaries. I'll do my best to prepare the most interesting materials from IT world. Let me know what you think! Let's begin, shall we?

## [Path Traversal in GitLab](https://bit.ly/3aNsohp) #security

Author of the vulnerability earned \$20,000 by finding a simple Path Traversal. The attack is trivial:

-   Create two projects
-   Add an issue with the following description:

```
![a](/uploads/11111111111111111111111111111111/../../../../../../../../../../../../../../etc/passwd)
```

-   Move the issue to the second project
-   The file will have been copied to the project

Old tricks still work, huh? :)

## [WebRTC in Go](https://bit.ly/2Yfob3h) #podcast #golang

This episode isn't only what WebRTC (spoiler alert - not only video!) is but also why Go was chosen of the project and why the author is happy about the choice.

## [Modlishka - reverse proxy in Go](https://bit.ly/3aNrbGK) #golang #network

Modlishka is open-source HTTP reverse proxy in Go. I wrote [a simple implementation of the reverse proxy](https://developer20.com/writing-proxy-in-go/) but this one is much more advanced. You can find a demo of how the proxy works.

## [Go reflections deep dive— from structs and interfaces](https://bit.ly/2xiBW69) #golang

Reflections are what I don’t like to do but sometimes there’s a need to play with it a bit. To avoid very stressful moments, I recommend reading about it before coding. This blog post is a good introduction to Go reflections.

## [Building a JSON Parser and Query Tool with Go](https://bit.ly/2zAKw0N) #golang

The interesting blog post where the author describes how he built a JSON query builder in Go. Building a custom lexer can be tricky so recommend to read it :)

## [Debugging with Delve](https://bit.ly/2Si1szX) #golang #debugging #delve

Still using `fmt.Print`, don't you? I do it as well (from time to time). Delve is an excellent tool which understands Go better than GDB. In more complex project, it's definiately worth trying,

## [How to Manage Database Timeouts and Cancellations in Go](https://bit.ly/2yQV63t) #golang #database

The `context.Context` is one of the most important structs in Go IMO. This blog post describes how to use it to manage cancellations and timeouts in Go with DB queries.

## [How To Make Life Easier When Using Git](https://bit.ly/3f7eavd) #git

Yet, another post about Git. What I learned from this post are commit templates. Maybe you’ll find something interesting as well :)

## ['Witcher 3' on the Nintendo Switch: CPU & Memory Optimization](https://bit.ly/3d0fqxY) #gamedev #performance

An awesome presentation about who Witcher 3 was ported to Nintendo Switch. Full of detailed information and tables :)

I hope you like the newsletter. I’ll try to track which topics are the most interesting for you to prepare better content. See you next week. Cheers!
