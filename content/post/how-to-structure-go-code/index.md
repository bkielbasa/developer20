---
title: "How to structure Go code?"
publishdate: 2021-10-27
categories:
    - Golang
    - Programming
tags:
    - golang
resources:
    - name: header
    - src: featured.jpg
---

> Programs should be written for people to read, and only incidentally for machines to execute - Abelson and Sussman

It is one of the most popular questions. You can find on the Internet attempts to answer this question. I've had concerns if I'm designing my packages or even the whole project correctly. Today, I'm not 100% sure about that! 

Some time ago, I had the pleasure to meet [Robert Griesemer](https://en.wikipedia.org/wiki/Robert_Griesemer) (one of Go's authors) in person. We asked him this question: "How to structure Go code?". He said: "I don't know." - It's not satisfying, I know. When we asked him about how he designs his code Robert said that he always starts with a flat structure and creates packages when he has to. That's much more concrete.

I've spent a lot of time trying different approaches in both pet projects in production applications. In this article, I'll show you options and tell you about the pros and cons of all of them. After reading this blog post, you won't have one "pattern to rule them all". You'll have a pocket knife in your toolbox.

## Before we start

Regardless of how you'll structure the code you have to think about everyone who reads it. The most important thing is that you shouldn't make your contributors or colleagues think. Put everything into obvious places. Don't reinvent the wheel. Do you remember what [Rob Pike](https://www.youtube.com/watch?v=PAAkCSZUG1c) said about `go fmt`?

> Gofmt's style is nobody's favourite, but gofmt is everybody's favourite.

You may not like well-known patterns. It's better for us, and for the whole community, to stick to them. However, if you have a good reason to make different decisions it's OK too. If your packages are well-designed the source tree will reflect it.

Let's start with the documentation. Every open-source Go project has his documentation on [pkg.go.dev](https://pkg.go.dev/). The very first thing you can see for every package is the overview of the package. Having `net/http` package as an example, before describing every public function, constants or types you have a description of what the package provides. You can learn from it how to use the API and more in-depth details. From which source the overview is generated? The `net/http` package has a [doc.go](https://github.com/golang/go/blob/master/src/net/http/doc.go) file where the author puts the general description of the package. You can put this overview to any other file in the folder but the `doc.go` is the standard one.

What should be in the `Readme` file then? Firstly, a general overview of this project — its goal. Then, you should have a quick start section where you describe what you should do to start working on the project, including any dependencies like Docker or any other tool we're using. You can find basic examples here or links to an external website where the project is described in more detail. You have to remember that what should be kept here is project-dependent. You have to think from our readers' point of view. What information is the most important for them?

When you have more documentation to provide put them into the `docs` folder. Don't hide them into folders like `/common/docs`. This approach has benefits: it's easy to find and in one pull request you can change both the public API and its documentation. You don't have to clone another repository and synchronize changes between them.

My next suggestion may be suppressing. Use well-known tools like `make`. I know that there are alternatives like [sake](http://tonyfischetti.github.io/sake/), [mage](https://github.com/magefile/mage), [zim](https://github.com/fugue/zim/) or [opctl](https://opctl.io/). The problem is that to start using them you have to learn them. If any project will be using a different automation tool it will be harder to get started for new maintainers. When you use the same (not that popular tool like `make`) in all our projects then that's fine. My point is that you should choose your tools wisely. Perfectly, our toolset should be obvious and boring. The more tools you use the more difficult it's becoming to start working on the project.

I've been working on a project where I had to run 2 different dependencies locally, be logged in to an AWS account in the CLI, and be connected to VPN to be able to run tests on my PC. The basic setup took a day or two to complete and I think I don't have to tell you how [flaky](https://testing.googleblog.com/2020/12/test-flakiness-one-of-main-challenges.html) those tests are.

For linting use [golangci-lint](https://github.com/golangci/golangci-lint). Enable all linters that seem to be reasonable for your project. In general, linters enabled by default may be a good start. 

## The flat structure (single package)

Let's start with the most recommended approach: keep the whole code in the root folder as long as you are not forced to add a new package. It is true, at the beginning of the project. I found it helpful when I got started and had a blurred idea of how it will work in the end.

Keeping everything in one place helps to avoid the cyclomatic dependency between packages. When you put something into a separate package and you'll find out that you need something else from the root folder, you'll be forced into creating a third package for this shared dependency. The situation is getting worse when the project evolves. You can end up having a lot of packages where most of them have dependencies on almost every other one. A lot of functions or types will have to be public. This situation blurs the API and makes it harder to read, understand and use.

With a single package, you don't have to skip between folders and think about the architecture because everything is in one place. It doesn't mean you have to keep everything in a single file.

```
courses/
  main.go
  server.go
  user_profile.go
  lesson.go
  course.go
```

In the example above, every logical part is organized into separate files. When you make a mistake and put a struct into the wrong file all you have to do is cutting and pasting it into a new location. You can think about it this way: a single file represents one solid part of your application. You can group your code by what it is (HTTP handlers, database repositories) or what it provides (managing user's profile). When you'll need something you'll know where you can find it.

When to create a new package? If you have a **good** reason to do so. What can it be?

1. When you have more than one way of starting your application.

Let's say you have a project and you want to run it in two modes: CLI command and web API. In this case, it's common to create a `/cmd` package with `cli` and `web` sub-packages.

```
courses/
  cmd/
    cli/
      main.go
      config.go
    web/
      main.go
      server.go
      config.go
  user_profile.go
  lesson.go
  course.go
```

You can put multiple `main()` functions into separate files into a single folder. To run them you have to provide an explicit list of files to compile with only one `main()` in it. This complicates running the application a lot. It's simpler to just type `go run ./cmd/cli`.

The usage of the `./cmd/` folder may sound over complicated when you have a single sub-package in it. I found it useful when needed to add, for example, consuming messages from a message broker. This topic will be covered in more detail in the section that [focused on splitting dependencies](#splitting-dependencies).

2. When you want to extract more detailed implementation

The standard library is an excellent example. Let's consider the [net/http/pprof](https://pkg.go.dev/net/http/pprof) package. The `net` package provides a portable interface for network I/O, including TCP/IP, UDP, domain name resolution, and Unix domain sockets. You can build any protocol you want based on what this package provides. The `net/http` gives us the ability to send or receive HTTP requests. The HTTP protocol uses TCP/UDP so it's natural that the `http` package is a sub-package for the `net` package. All types and methods in `net/http/pprof` package are available throw the HTTP protocol so the natural consequence is a child package for the `http` one.

The same is true for `database/sql` package. If you'll have more implementation for non-relational databases they will go under `database` package, next to the `sql` package.

Can you see the pattern? The deeper the packet is in the tree, the more details deliver. In other words, every sub-package provides a more concrete implementation of the thing on the parent package.

3. When you started to add a common prefix to closely related things

After some time, you may notice that to avoid misunderstandings or naming conflicts, you started adding a prefix or suffix to functions or types. It may be a good sign that by doing it we're trying to simulate the lack of packages in our project. It's hard to say when it's time to extract the new sub-package. Do it every time when you'll see that it improves the API's readability and make the code clearer.

```go
// TODO: think about a better example!
r := networkReader{}

//

r := network.Reader{}

```

As you can see, the flat structure is simple and powerful at the same time. There are use cases where you may find it useful and helpful. This way of organizing the code works not only for small or greenfield projects. Here are examples of libraries that follow the single package pattern:

* https://github.com/go-yaml/yaml
* https://github.com/tidwall/gjson

It's worth remembering that you shouldn't stick to this way of organising the code at all costs. There are reasons for keeping it simple but adding a few more packages may make your code better. Unfortunately, there's no one silver bullet. What you have to do is experimenting and asking your colleagues or maintainers which option is more readable for them.

## Modularisation

The previously described way of organising your code may not be efficient in use cases. I've spent a lot of time trying to get 'the right' project structure. After a time, I noticed that for types of applications I started organising the code in a similar way. It was business applications.

When we're working on applications that directly give our clients values, the flat structure maybe not be efficient. You want to create modules that provide a set of features related to controllers, infrastructure, or part of the business domain. Let's take a closer look at two of the most popular ways of doing it and talk about their pros and cons.

### Organising by kind

This model is popular. Nobody I know is advocating using this strategy for organising code but I find it in both old and new projects. Organisation by kind is a strategy that tries to bring order to overly complex units of code by throwing the parts into buckets based on which structure it is. It's common to have packages called `repositories` or `model`. A consequence of doing is creating packages like `utils` or `helpers` where you have a feeling that you should put a function or a struct to a separate place but there's no good spot for it anywhere.

```
.
├── handlers
│   ├── course.go
│   ├── lecture.go
│   ├── profile.go
│   └── user.go
├── main.go
├── models
│   ├── course.go
│   ├── lecture.go
│   └── user.go
├── repositories
│   ├── course.go
│   ├── lecture.go
│   └── user.go
├── services
│   ├── course.go
│   └── user.go
└── utils
    └── stings.go
```

In the example above, you can see that project is organised by type. When do you want to add a new feature or fix a bug related to a course, where would you start looking? At the end of the day, you'll start jumping from a package to a package hoping you'll find something useful there.

![Graph that shows dependencies between packages](./modularization.png)

This approach has its consequences. Every type, constant, or function has to be public to be accessible in another part of the project. You'll end up with most types marked as public. Even for those that shouldn't be public. It makes confusion about what's important in this part of the application. Many of them are details that may change at any time.

On the other hand, organizing by kind is natural for us. We are technical people that think in categories of handlers or database tables. This is how we grew up and how we have been taught. If you're not experienced, this approach may be more beneficial because it may help you get started faster. In a long run, you may experience inconveniences but it doesn't mean your project will fail - the contrary. There's a lot of successful applications that are designed this way.

### Organising by components

A component is part of the application that provides an independent feature that has little or no dependencies outside. You can think about it as plugins that when you plug out one of them, the whole application works but has limited functionality. It may happen in production applications that are running for months or years.

The application may have one or more core components that deliver the business values. In the Domain-Driven Design terminology the component is a bounded context. We will describe  the DDD with the context of Go in another chapter. TODO: ADD THE REFERENCE.

The package's API should describe what the package provides and not more. It shouldn't expose any low-level details that aren't important from the consumer's point of view. It should be as minimalistic as it's possible. The consumer may be another package or another developer who imports our code.

The component should contain everything it needs to provide the business value. It means, every storage, HTTP handler, or business model should be stored inside of the folder.

```
.
├── course
│   ├── httphandler.go
│   ├── model.go
│   ├── repository.go
│   └── service.go
├── main.go
└── profile
    ├── httphandler.go
    ├── model.go
    ├── repository.go
    └── service.go

```

Thanks to organizing the code this way, when you have a task-related course, you know where to start looking. It's not spread across the whole application. However, it's not easy to achieve good modularisation. It may require multiple iterations to achieve a good package API.

There's one more challenge. What if these packages depend on each other? Let's say that you want to display the most recent courses on the user's profile. Should they share the same repository or a service?

In this particular situation, from the profile's point of view, courses are an external dependency. The best way of tackling the problem is to create an interface in the `profile` package that has the required method.

```go
type Courses interface {
  MostRecent(ctx context.Context, userID string, max int) ([]course.Model, error)
} 
```

In the `course` package you expose a service that implements this interface.

```go
type Courses struct {
  // maybe some unexported fields
}

Func (c Courses) MostRecent(ctx context.Context, userID string, max int) ([]Model, error) {
  // return most recent coursers for specific user
}
```

In the `main.go` you create the instance of the `Courses` struct from the `course` package and pass it to the `profile` package. In the test in the `profile` package, you create a mock implementation. Thanks to this, you can develop and test the profile functionality without even having a `course` package implemented.

As you can see, the modularisation makes the code more maintainable and readable but it makes you think harder about your decisions and dependencies. The logic may look like a perfect fit for a new package but seems to be too small. On the other hand, during working on the project part of the existing package may start growing and, after time, promoted to an autonomous piece of code.

When the code grows inside of the package, you may ask yourself: how to organize the code inside of the single module? That's another hard question to answer. In this section, I showed the flat structure while using components of the application. But, sometimes it's not enough...

## Clean Architecture

You probably heard the term: [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html), Onion Architecture, or similar. Uncle Bob [wrote a book](https://www.amazon.com/Clean-Architecture-Craftsmans-Software-Structure/dp/0134494164) that describes in detail what every layer means and what should or shouldn't contain. The idea is simple. You have 4 layers of your application or module (depending on how big your codebase is): Domain, Application, Ports, Adapters. In some sources, names may differ. For example, instead of using Ports and Adapters, the author uses Inbound and Outbound. The core idea says the same. Let's describe every layer with examples.

### Domain

This is the heart of our application. Every business logic should live here. It means if any business requirement is changed or added, you have to update our Domain part. This package should have **no** dependencies on anything outside. It should not *know* in which context this code is executed. It means, it should not rely on any infrastructure part or know any UI details.

```go
course := NewCourse("How to use Go with smart devices?")
s := course.AddSection("Getting started")
l := s.AddLecture("Installing Go")
l.AddAttachement("https://attachement.com/download")
// etc
```

Please notice that at this point you don't care about where to store the course or how you can add a new one (using an HTTP request or using a CLI). In the `domain` package you describe what the course may contain and what operations you can do on it. That's all.

### Application

This layer holds every use case of your application. It's a glue point between infrastructure and the domain. In this place, you get the input (from wherever it comes from), apply it on the domain objects, and persist or send it somewhere else.

```go
func (c Course) Enroll(ctx context.Context, courseID, userID string) error {
  course, err := c.courseStorage.FindCourse(ctx, courseID)
  if err != nil {
	  return fmt.Errorf("cannot find the course: %w")
  }
  
  user, err := c.userStorage.Find(ctx, userID)
  if err != nil {
	  return fmt.Errorf("cannot find the user: %w")
  }
  
  if err = user.EnrollCourse(course); err != nil {
	  return fmt.Errorf("cannot enroll the course: %w")
  }
  
  if err = c.userStorage(ctx, user); err != nil {
	  return fmt.Errorf("cannot save the user: %w")
  }
  
  return nil
} 

```

In the code above, you can find a use case where a user enrolls in a course. It is the combination of two pieces: interacting with the domain objects (User, Course) and infrastructure (storing and fetching the data).

### Adapters

Adapters are also known as Outbound or Infrastructure. This layer is responsible for storing and fetching data with the outside world. It can be a database, blob storage, a file system, or another (micro) service. Very often this layer has its representation in interfaces in the application layer. It helps to test the application layer without running the database or writing files to the file system.

The Adapter is an abstraction over a low-level detail so other parts of your software don't have to "know" which database version you're using, how the SQL query looks like, or where you store your files.

### Ports

Ports (known as Inbound) are this part of the application that's responsible for getting the data from the user. It can be an HTTP handler, an event handler, or a CLI command. It gets the user's input and passes it to the application layer. The outcome of this operation goes back to the Port.

```go
func enrollCourse(w http.ResponseWriter, r *http.Request) {
body, err := io.ReadAll(r.Body)
	if err != nil {
	  w.WriteHeader(http.StatusBadRequest)
	  logger.Errorf("cannot read the body: %s", err)
	  return
	}
	
	req := enrollCourseRequest{}
	if err = json.Unmarshal(body, &req); err != nil {
	  w.WriteHeader(http.StatusBadRequest)
	  logger.Errorf("cannot unmarshal the request: %s", err)
	  return
	}
	
	if err = validate.Struct(req); err != nil {
	  w.WriteHeader(http.StatusBadRequest)
	  logger.Errorf("cannot validate the request: %s", err)
	  return
	}
	
	if err = app.EnrollCourse(req.CourseID, req.UserID); err != nil {
	  w.WriteHeader(http.StatusInternalServerError)
	  logger.Errorf("cannot enroll the course: %s", err)
	  return
	}
}
```

Please notice that writing a CLI command that executes the same logic is straightforward. The only difference is the source of the input.

```go
var userID string
var courseID string

var enroleCourseCmd = &cobra.Command{
	Use:   "courseID userID",
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err = app.EnrollCourse(courseID, userID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Errorf("cannot enroll the course: %s", err)
			return
		}
	},
}
```

Keeping these layers clean and consistent may bring a lot of value to your code. It's easy to test, responsibilities are clear and it's much more obvious where to start looking for the code you want to change. If it's a bug related to courses and it's a business logical issue, you'll start checking the domain or application layers.

On the other hand, it's difficult to keep the boundaries clear and consistent. It requires a lot of self-discipline, experience, and at least a few iterations to do it right. That's why many people fail in this field.

## Summary

Structuring the code is hard. What makes it even more difficult is the fact that the architecture of the application may change during its lifetime a few times. It evolves. You may start with a flat structure but end up with multiple modules with many sub-packages. Don't expect to get it right the first time. It may require multiple iterations and gathering feedback from others.

What's more, you may mix different ways of organizing the code depending on the part of the application. In a place where your business logic lives, you will start with modularisation. However, many applications need utilities that don't fit into any of the existing packages. You may follow the flat structure pattern there.

