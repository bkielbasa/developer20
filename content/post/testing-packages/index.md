---
title: Testing packages APIs
publishdate: 2021-12-16
categories:
    - Golang
    - Programming
tags:
    - testing 
    - golang
resources:
    - name: header
    - src: featured.jpg
---

There are many practices and tactics that tackle testing. Today, I'll share with you how I write tests in my projects. Please notice that you may find it useful when you're starting a new project or an independent part of existing applications. You may find it difficult to apply in an already existing application. It's not impossible but may be challenging.

## How does the architecture of the package look like?

Before I show you my approach, I have to explain how I design packages. You can read in more detail in another [blog post](https://developer20.com/how-to-structure-go-code/). Today, I'm focusing on the "Clean Architecture", one that I found the most useful in business-focused applications.

I like splitting tests into three parts: [given, when and then](https://martinfowler.com/bliki/GivenWhenThen.html). In the `given` section I prepare everything. In the `when` section (often it's 1 line long) I operate I want to test. In the last `then` part, I make assertions. You can see an example below.


```go
var storage app.ProductStorage

func TestFetchingProductInTheCatalog(t *testing.T) {
	is := is.New(t)
	// given
	ctx := context.Background()
	appServ := app.NewProductService(storage)

	productID, err := addNewProduct(ctx, storage)
	is.NoErr(err)

	// when
	fetched, err := appServ.Find(ctx, productID)

	// then
	is.NoErr(err)
	is.NoErr(productEquals(p, fetched))
}
```

Let's ignore the `storage` variable for a while. While designing both tests and the API I think about how many things I can hide to make the test more and more readable. That's why I created the `buildProduct()` methods that just build the product. Details about how exactly it looks like aren't important here so I extracted it to the function. We can go even further and extract adding the product to the storage.

```go
	productID, err := addNewProduct(ctx, storage)
	is.NoErr(err)
```

Here comes the interesting part - the `storage` variable. As you can see, its type is `app.ProductStorage`. It's an interface. I prepared two implementations of the interface. The first one is in-memory (used mostly in tests). Initializing the in-memory version is straightforward.

```go
//go:build !integration

package tests

import "github.com/bkielbasa/go-ecommerce/backend/productcatalog/adapter"

func init() {
	storage = adapter.NewInMemory()
}
```

Please notice the build tag I added at the top of the file. It says: "compile this file as long as the `integration` build tag isn't provided". All tests, by default, should be fast and reliable. We'll want to run them very frequently so we shouldn't wait for their results too long. What's more, it's a good practice when someone downloads our project and type `go test ./...` all tests are passing.

There's nothing more annoying than reading docs and setting up everything just to be able to run tests. Making this first experience a pleasure improves morale.

What if we want to test against a real database? I create a separate file (I put it next to the previous one) with connecting to everything I need.

```go
//go:build integration

package tests

import "fmt"
import "github.com/bkielbasa/go-ecommerce/backend/productcatalog/adapter"
import "os"

func init() {
	pass := getEnv("POSTGRES_PASSWORD", "")
	var conn string

	if pass != "" {
		conn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", getEnv("POSTGRES_HOST", "localhost"), getEnv("POSTGRES_PORT", "5432"), getEnv("POSTGRES_USER", "bartlomiejklimczak"), pass, getEnv("POSTGRES_DB", "ecommerce"))
	} else {
		conn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", getEnv("POSTGRES_HOST", "localhost"), getEnv("POSTGRES_PORT", "5432"), getEnv("POSTGRES_USER", "bartlomiejklimczak"), getEnv("POSTGRES_DB", "ecommerce"))
	}

	s, err := adapter.NewPostgres(conn)
	if err != nil {
		panic("cannot establish connection to postgres: " + err.Error())
	}

	storage = s
}
```

I get Postgres credentials from env variables and initialize the Postgres adapter. This test expects that all migrations were executed and the application has access to the database. On the top of the file, I added `go:build integration` build tag that says: "if the `integration` build tag is provided, compile this file".

Right now, when I type `go test ./...` when I want quick feedback. To run **the same** tests but with a real database or message broker I type `go test ./... -tags integration`. This approach has one big advantage.

When I want to change something in the logic of the application I can run quick tests only. I run integration tests only before pushing to a remote branch. Both tests should be run in CI/CD as well.

![tests result in github actions](./tests.png)


Please notice that this approach works mostly on business-focused applications. If you're writing a library or a tool, the logic inside of your code may not be that significant as it is in business applications.

## Summary
I showed you how it may work on package level but nothing stops you from writing similar tests but with end-to-end in mind. You can run the whole application and test your HTTP API with mocked or real databases under the hood. Thanks to this you'll write less number of tests that should cover most of the success paths.

The drawback is the need for implementing two implementations of the same thing. It may take some additional time. In my opinion, in long term, it pays off.

Please let me know what you think about it or share with you ways of writing tests in the comments section below.