---
title: "Testing packages"
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

There are many practices and tactics that tackle testing. Today, I'll share with you how I write tests in my projects. Please notice that you may find it useful when you're starting a new project or a independent part of existing applications. You may find it difficult to apply in already existing application. It's not impossible but may be challenging.

## How does the architecture of the package look like?

Before I show you my approach, I have to explain how I design packages. You can read in more details in another [blog post](https://developer20.com/how-to-structure-go-code/) about how various ways of solving the problem look like. Today, I'm focusing on the "Clean Architecture" one that I found the moste useful in business focused applications.



```go
var storage app.ProductStorage

func TestFetchingProductsInTheCatalog(t *testing.T) {
	is := is.New(t)
	// given
	ctx := context.Background()
	appServ := app.NewProductService(storage)

	p, err := buildProduct(ctx, storage)
	is.NoErr(err)
	err = storage.Add(ctx, p)
	is.NoErr(err)

	// when
	fetched, err := appServ.Find(ctx, string(p.ID()))

	// then
	is.NoErr(err)
	is.NoErr(productEquals(p, fetched))
}
```