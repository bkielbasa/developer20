---
title: "Go Programming Language - book review"
publishdate: 2021-02-10
resources:
    - name: header
    - src: featured.jpg
---

It took a while since I got this book. At the very beginning, I didn't want to buy it. However, I got so much feedback that this book is so good that I had to check it myself.

It is the second Go book in my library. I didn't write about the first one because I wanted to compare it with another one. I wanted to recommend you only the best. As you can see, I have a review about "The Go Programming Language". I'll go the book chapter by chapter and tell you what I like or don't like about each of them.

The first chapter is a tutorial - a huge plus from me for that. In other books, in general, the author tells about the syntax and gives simple examples. You'll learn how to work with command-line arguments or how to read files from the filesystem. You'll write a simple HTTP handler and understand some basics about how to work with requests and responses. Yeah, I like the practical way of learning.

The second chapter describes program structure like names, variables, type declarations, and scopes. I think that's what every good book or a tutorial about Go should have. The same counts in the next chapters. Chapter 3 describes basic data types like booleans, strings, integers, etc. What's more, it tells about binary operations that, in my opinion, are not needed at this point. If someone is starting to learn Go, a such thing can only confuse the reader. It's a good candidate for an appendix.

On the other hand, I like how the author described slices and arrays in chapter 4. The knowledge about how they work can save you from hard to find bugs. I even wrote an article "What you should know about Go slices". It's described clearly. Must-read for me.

The chapter about structs is OK. I mean, there's everything that's needed to know. You can find many other places where you'll find precisely the same knowledge. The author deserves praise for describing working with JSON format. This format is a standard in the IT world. The fact you can find it in the book is helpful.

There are other more designer-friendly ways of working with templates. In chapter 4 the authors tell you about the Go template engine. It's fine for a start, but you have to read much more about them to use them correctly (security and more). I don't know about you, but I'm not a fan of the syntax but it's just my personal objection.

The next chapter is about functions. And again - you can find almost everything that you need to know about them there. Kudos for writing about error handling in Go. Unfortunately, the book was published in 2016. It means the author doesn't describe the [latest changes in the language](https://blog.golang.org/go1.13-errors). After this chapter, I recommend reading the linked blog post.

The 6th part of the book describes methods. IMO this chapter is more important than the previous one because it describes the difference between [pointer and value semantics](https://developer20.com/pointer-and-value-semantics-in-go/) what can help you a lot at the beginning.

Interfaces are one of the most crucial features in Go. In my opinion, of course! You have to read the chapter more than once. Seriously. You'll write better code when you understand how Go interfaces work. What I like is that some examples use interfaces you'll use like `http.Handler`. When you start writing your web apps it will give you a feeling that you already know it. That's cool. You will find type assertions there as well as type switches and an example with parsing XML files.

There's one thing I think it's missing there. I didn't find a good enough description of how to design good interfaces. This is a problem in general, I know but this is the only thing I'd add to the chapter.

In the 8th and 9th chapters, you'll find a description of goroutines and channels as well as a description of `sync` package. Honestly? It's awesome! As long as you won't have problems with the Go scheduler, the chapter is all you need to know. You have examples there, a description of how channels work, cancellation, etc.

The 10th chapter is a bit out-of-date because we have Go modules and the GOPATH goes into oblivion. You can read it but after that, you should read the [Use Go Modules](https://blog.golang.org/using-go-modules) from the official blog.

In the chapter about tests, you'll learn how to write tests in Go! There are examples with table tests, white-box tests as well as generating a code coverage report. Basically, all you need to write tests in Go.

I have the biggest doubts about the last 2 parts of the book. The reflection and cgo is not what we try to use on daily basics. On the other hand, they can be handy in some cases... Writing your own `Sprintf` function can be a great adventure, too! In the beginning, you can skip those chapters.

My final word? If you have doubts if buy it - buy it. It's an excellent book that can help you understand better the Go programming language. I think that Alan Donovan and Brian Kernighan did a great job. There are a few things I'd change or add to it but after reading it you'll be ready to write your applications including CLI tools and web applications.

