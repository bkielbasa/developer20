---
title: Go web frameworks
publishdate: 2020-08-10
resources:
    - name: header
    - src: featured.jpg
categories:
    - Golang
tags:
    - golang
    - web framework

short: Go has plenty of different web frameworks. When you are faced with choosing a framework for the first time, it may turn out to be quite a challenge to choose the best one. This article is intended to help you choose the best one. It is full of personal judgments that you may disagree with. However, I believe you will find it most helpful.

---

Go has plenty of different web frameworks. When you are faced with choosing a framework for the first time, it may turn out to be quite a challenge to choose the best one. This article is intended to help you choose the best one. It is full of personal judgments that you may disagree with. However, I believe you will find it most helpful.

## [Martini](https://github.com/go-martini/martini)

The first framework is Martini. Honestly, it shouldn't be here as it's been under development since 2017. However, I added it because I found it on other compilations. Just don't use it :)

## [Gin](https://github.com/gin-gonic/gin)

Gin is probably the most popular Go web framework. They say they are the fastest as well. From the props is easy routing with named parameters. It's trivial when you want to get a form or query parameters. It speeds up the development. From things I really like is easy testing the HTTP part and very good documentation.

```go
func main() {
	router := gin.Default()

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
	router.Run(":8080")
}
```

What I don't like is that every handler has to implement their custom interface instead of the standard one.

```go
router.GET("/welcome", func(c *gin.Context) {
	// here's the code
})
```

Why do I care? If you have any middleware written for the standard HTTP handler or other Go framework, you have to rewrite it again. It's a huge vendor lock-in for me.

## [Buffalo](https://gobuffalo.io/en/)

I think that Buffalo has the best documentation. It supports hot reload which combined with frontend pipelines sounds very helpful. It uses [Plush](https://github.com/gobuffalo/plush) for templating. Be honest - the go build-in templates are good but aren't very user-friendly. Buffalo has support for named parameters as well but has the same problem with the signature of the HTTP handler function.

```go
func (c buffalo.Context) error {
  // do some work
}
```

The thing I love here is the CLI command. It gives a bit of [RoR](https://rubyonrails.org/) experience that's very nice. This framework looks like Go on Rails so if you like a similar approach then Buffalo is for you.

## [Goji](https://github.com/goji/goji)

The project isn't maintained anymore. Goji is a minimalistic and flexible HTTP request multiplexer for Go. As you may expect - you won't find many features here. Its docs are very minimalistic.

## [Gorilla](https://www.gorillatoolkit.org/)

Gorilla is a toolkit with a number of packages available. You can find the [mux](https://github.com/gorilla/mux) for HTTP routing, [csrf](https://github.com/gorilla/csrf) for Cross Site Request Forgery prevention or [feed](https://github.com/gorilla/feeds) for rss/atom generation. The downside of it it doesn't give you a solid working framework but a list of libraries you can use to build your own framework.

As you can see, there are several options to choose from. It's hard to choose the best because there's no best Go framework. It depends on your needs and preferences. However, I hope I gave you a clue about what options we have and recommend you to give a try at least 2 projects.
Do you use any other frameworks or have any thoughts? Leave them in the comments section below. 