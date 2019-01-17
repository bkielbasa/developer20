---
layout: post
title:  How to name exceptions? It’s not so obvious…
mainPhoto: exceptions.jpeg
---

Naming things is one of the most difficult things in our job. Naming exceptions are even more complicated because exceptions are not regular classes. In this article, I’ll tell you a bit about naming conventions I’ve found and tell pros and cons of them.
## The ‘Exception’ suffix. To add or not to add

There are two schools to add the suffix or not to add it. Let’s take a look a bit closer to see main pros and cons of both points of view.
### Leaving without the suffix

Exceptions are not a regular class. In high-level languages like Java, C#, Python or PHP, you cannot throw or catch something that’s not an exception. If something is thrown, you know it’s an implementation of a specific interface. Adding an extra word, in the end, is just redundant. It doesn’t give you any useful information you did not already have. It only makes the name of the class longer. If you’re dealing with an exception and don’t recognize it as such, usually something is deeply wrong. Moreover, every time you’re facing with exceptions, you have some keywords like try, catch, finally which say that.

Of course, in languages like c++ you can throw anything. For example, you can throw a null (BTW, it will be cast to long in this case).
### The world with the ‘Exception’ suffix

Imagine you’re implementing a parser and you have a completely valid classMissingParameter. You use it to display some debug messages in a console. In another place, you have an exception which is thrown when someone did not put a required parameter. Of course, its name is MissingParameter. Many IDE’s let you browse and search classes in a project.

It may be confusing. You have two classes with the same name (in different namespaces, of course) but one of them is an exception, the second isn’t. I know you can see the folder where the class is placed. If it’s in folder `exception` it’s probably an exception. The suffix will make it more clear. In a small project it’s not so huge problem but in a large one, it might be.
## The naming

What does say `InvalidArgumentException` you? An invalid argument… what? What did just happen? Did a user not enter a required parameter or he did it but put a null or put some other wrong value? Maybe the class passed as the argument has a wrong state? When you have only information from the exception you do not know exactly what’s going on.

It’s a reason why the exception should be a noun and should describe exactly what happened. A more clear exception would be `ArgumentCannotBeNull` or `ProductIsOutOfStock`. Thanks to it when you take a look at the `try-catchstatement` you know exactly what errors may occur.

On the other hand, you’ll have to create a new class for every unwanted behavior. It can be really a large number of new classes! Imagine that a class can throw 10 exceptions. How will the catch section look like? In many cases, you’ll deal with the exception the same way. It’s not only redundant but you have to duplicate your code.

To solve it you can create a more generic class like `InvalidArgumentException` with a custom message. The message will tell you more about the problem and it can be more dynamic for example ‘A product with SKU “how-to-handle-fame” is out of stock’. It’s much handier than an exception `ProductIsOutOfStock(“how-to-handle-fame”)`.
## Summary

As I said in the introduction, naming things is one of the most difficult tasks we have as developers. Good class’ name or variable’s name may save hours of debugging or just reading the code you try to understand. In many cases, it isn’t obvious or clear in the beginning. Your decisions may impact your code in many ways so do not be afraid of chasing your mind and refactoring code.

Please let me know if you have any other thoughts about exceptions, please let us know in the comments section.