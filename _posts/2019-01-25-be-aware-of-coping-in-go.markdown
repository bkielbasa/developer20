---
layout: post
title:  Be aware of copying in Go
mainPhoto: copy-in-go.jpg
categories: Golang
tags: [golang, channels]
---

Some bugs are very hard to find and to reproduce but easy to fix. To avoid them, it's helpful to know how the tools we're using work under the hood. From this article, you'll learn what shallow and deep copy are and which errors you can avoid thank's the knowledge about them.

Can you find a problem with the code below?

```golang
q1 := NewQuestion(1, "How to be cool?")
q1.AddAnswer(1, "eat peanuts")

q2 := q1

q2.AddAnswer(2, "visit developer20.com regulary!")
fmt.Println("How to be cool?")
q1.ShowAnswers()
```

It's hard to say what's wrong because, at first glance, the code seems to work OK. When we run the code we receive the output:

```
How to be cool?
 * eat peanuts
 * visit developer20.com regulary!
```

The second answers were added to the first question unexpectedly. To find out where's the problem, let's take a look at the struct's definition.

```go
type Question struct {
	ID int
	Content string
	Answers map[int]string
}

func NewQuestion(id int, content string) *Question {
	return &Question{
		ID: id,
		Content: content,
		Answers: make(map[int]string),
	}
}

func (q *Question)AddAnswer(id int, content string) {
	q.Answers[id] = content
}

func (q *Question)ShowAnswers() {
	for _, answer := range q.Answers {
		fmt.Printf("* %s\n", answer)
	}
}
```

[Golang play](https://play.golang.org/p/X9T_EGSJ7Hk)

The problem is in `Answers` definition in `Question` structure. In the documentation, you'll find that

> Slices, maps and channels are reference types that do not require the extra indirection of an allocation with new.

It means that the values of the array are kept in different place in the memory but in the struct only a reference to them is kept. By simply copying the struct, the reference is copied. The copied reference points to the same address in the memory.

![](/assets/posts/struct-copy2.png)

This is an example of a shallow copy. Shallow copy happens when an object is copied byte by byte. If the object which is copied contains a reference (or a pointer) its address is copied. This situation is illustrated in the above picture. To avoid such situation, deep coping has to be implemented manually.

Why is it important? Imagine a situation where a simple struct with only basic types. The coping of the struct is safe. Hoverer, after some time a referenced type field was added. It can be a slice or a map. It is possible that tests won't cover this edge case. Sometimes a bug report is received from production users about some situations where the system works unpredictably. Does this scenario sound possible? [But it happens](https://allegro.tech/2017/07/golang-slices-gotcha.html). This bug is very similar to the above one but is related to slices. 

## Summary

Even simple copying can lead to serious bugs which are extremely hard to find. In Go, there's no way to prevent from coping a struct but in c++ you can without any problem.

```c++
SomeClass(const SomeClass&) = delete;
```

In Java, the class has to implement `Cloneable` interface to be cloneable. To prevent those problems, it's better to implement own function for the cloning purpose or cover your code with good tests.

I hope you enjoyed the article. Any thoughts? Do you have anything else to add? Feel free to leave your comment below. Cheers!