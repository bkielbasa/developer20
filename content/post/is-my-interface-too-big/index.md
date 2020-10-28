---
title: "Is my interface too big?"
publishdate: 2020-10-28
resources:
    - name: header
    - src: featured.jpg
categories:
    - Golang
tags:
    - golang
    - code-review
---

In this article, I explain how you can detect if the interface you're using is getting too big and requires splitting into smaller ones. Smaller interfaces help to improve the maintenance and readability of the code. What's more, it helps with understanding the code.

Interfaces in Go are different than those known in Java, c#, PHP etc. In those languages you define interfaces up-front. In other words, at the moment of creating a class you need to know how the class will be used. In Go things are different. You can create a struct with as many functions you want and the user of if can define only a sublist of methods he needs. It's very powerful tool. But sometimes we still can create too big interfaces. The list above should help you find interface segragation issues in your code.

## Panics in mocks

If you're generating/writing mocks and some of the methods have empty implementation or a `panic` inside, you're probably incorrectly segregated interfaces. This code is a good suggestion that you have more than one responsibility in the code and you're trying to test one of them.

My suggestion is to try to segregate those responsibilities into separate functions/structs and test them separately. Here's an example.

{{< highlight go >}}
type MethodRepoMock struct {
  mock.Mock
}

func (m *MethodRepoMock) Create(ctx context.Context, method dto.ShippingMethodDTO, tx *sql.Tx) error {
  panic("implement me")
}

func (m *MethodRepoMock) GetAll(ctx context.Context, sellerID string, templateID string, countryCode *string) ([]*models.ShippingMethodAggregate, error) {
  panic("implement me")
}

func (m *MethodRepoMock) Get(ctx context.Context, sellerID string, methodID string) (*models.ShippingMethodAggregate, error) {
  args := m.Called(sellerID, methodID)
  if args.Get(0) == nil {
    return nil, args.Error(1)
  }
  x, ok := args.Get(0).(*models.ShippingMethodAggregate)
  if !ok {
    panic(fmt.Sprintf("assert: arguments failed because object wasn't correct type: %v", args.Get(0)))
  }
  return x, args.Error(1)
}

func (m *MethodRepoMock) Update(ctx context.Context, method dto.ShippingMethodDTO, sellerID string, methodID string, tx *sql.Tx) error {
  panic("implement me")
}
{{< / highlight >}}

In the code above, you can clearly see that only one function is used in the test case. The service that uses the mock can be updated to use a simple and small interface for the repository.

{{< highlight go >}}
type Getter interface {
  Get(ctx context.Context, sellerID string, methodID string) (*models.ShippingMethodAggregate, error)
}
{{< / highlight >}}

## Changes in not related code

Another premise saying that your interface is too large is a situation when you add a new method to it and it requires some changes in other areas of the code - not related to the change you're making.

Let's say you have two services (`Service1` and `Service2`) and a repository. The repository is used in two different services (what's not a bad thing). A new requirement came and you have to add a new function in `Service1` and into the repository. If the change requires changes, for example, in tests for the `Service2` it a sign that there's bad interface segregation.

This problem can happen in two scenarios: both services use a concrete repository or share the same (larger) interface.

{{< highlight go >}}
type myRepository interface {
  // used only in Service1
  GetSomething(ctx context.Context, id string) (Something, error)

  // used only in Service2
  CalculateSomething(ctx context.Context, param1 int) (int, error)

  // used only in Service2
  SendSomething(ctx context.Context, param1 int) (int, error)

  // used in both services
  SaveSomething(ctx context.Context, id int) (int, error)
}

type Service1 struct {
  repo myRepository
}

type Service2 struct {
  repo myRepository
}
{{< / highlight >}}

Don't be afraid of creating small interfaces, even if some method in the interface can repeat. When we refactor the code above we'll end up with two interfaces that have the same functions in it.

{{< highlight go >}}
type serice1Repo interface {
  // used only in Service1
  GetSomething(ctx context.Context, id string) (Something, error)

  // used in both services
  SaveSomething(ctx context.Context, id int) (int, error)
}

type serice2Repo interface {
  // used only in Service2
  CalculateSomething(ctx context.Context, param1 int) (int, error)

  // used only in Service2
  SendSomething(ctx context.Context, param1 int) (int, error)

  // used in both services
  SaveSomething(ctx context.Context, id int) (int, error)
}
{{< / highlight >}}

This kind of duplication is fine. We have the tendency to follow the [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself) principle too aggressively and create an imagined being just to save a few lines of code.

{{< highlight go >}}
type servicesRepo interface {
  // used in both services
  SaveSomething(ctx context.Context, id int) (int, error)
}

type serice1Repo interface {
  servicesRepo
  // used only in Service1
  GetSomething(ctx context.Context, id string) (Something, error)
}

type serice2Repo interface {
  servicesRepo

  // used only in Service2
  CalculateSomething(ctx context.Context, param1 int) (int, error)

  // used only in Service2
  SendSomething(ctx context.Context, param1 int) (int, error)
}
{{< / highlight >}}

As you can see, you don't even save lines of code. You have them even more! IMO it doesn't improve the readability as well.

## Public interface

Making a public interface isn't a bad thing, though. The reason why you're doing it can be. Let's go back to our two services and the repository. If you have to make the repository public because those services have to have access to them (and both of them are in different packages) it's a code smell and a signal that it needs refactoring.

I've seen many times a package called `repositories` where all interfaces have their place. Every service that wanted to use one of the interfaces has to take the whole (bigger) interface and, as we saw in the "Panics in mocks" section, ignore other unused methods. How to solve the problem?

Go's interfaces are awesome because you can fix it with low effort. The only thing you have to do is to put small interfaces with only methods you need next to the code that uses it and use it instead. When you'll keep refactoring it step-by-step you'll hit the point where you will be able to remove the `repositories` package completely.

## Summary

That's all I have for you this time. I hope you found this article interesting and I helped you build better software. Those suggestions aren't Go-specific but I put them in the context to help understand the principles. I believe there are more signals of bad interface segregation. If you have your own idea - let us know in the comments section below.

