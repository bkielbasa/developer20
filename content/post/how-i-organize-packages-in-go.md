---
title: How I organize packages in Go
publishDate: 2019-08-22
categories: [Golang, Programming]
tags: [golang]

---

Structuring the source code can be as challenging as writing it. There are many approaches to do so. Bad decisions can be painful and refactoring can be very time-consuming. On the other hand, it’s almost impossible to perfectly design your application at the beginning. What’s more, some solutions may work at some application’s size and should the application [develop over time](https://www.amazon.com/Building-Evolutionary-Architectures-Support-Constant-ebook/dp/B075RR1XVG/ref=sr_1_1?keywords=evolutionary+architecture&qid=1565498731&s=gateway&sr=8-1). Our software should grow with the problem it’s solving.

I mostly develop microservices and this architecture fits great to my needs. Projects with much more domain in it or more infrastructure applications may require a different approach. Let me know in the comments below what’s your design and where it makes sense the most.

## Packages and its dependencies

When it comes to developing domain services, it’s useful to split the service by domain’s components. Every component should be independent and, in theory, be able to be extracted to an external service if needed. What does it mean and how to achieve that?

Let’s assume that we have a service which handles everything related to placing orders like sending an email confirmation, saving information to a database, connecting with a payment provider etc. Every of the package should have a name which clearly [says what’s for](https://www.amazon.com/Clean-Architecture-Craftsmans-Software-Structure/dp/0134494164) and compatible with [the naming standard](https://blog.golang.org/package-names).

![](/images/organize-go.png)


This is only an example of a project where we have 3 packages: `confemails`, `payproviders` and `warehouse`. Names should be short and self-explaining.

Every of the package has his own Setup() function. The function accepts only bare requirements to be able to work correctly and be able to communicate with the world outside of the package. For example, if the package exposes an HTTP endpoint, the Setup() function accepts an HTTP router like mux Router. When you’re package requires access to the database then Setup() function accepts sql.DB as well. Of course, the package can need another package, too.
