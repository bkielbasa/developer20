---
title: GoBDD
publishDate: 2020-12-28

---

[GoBDD](https://go-bdd.github.io/gobdd/) is a library that helps writing BDD tests in Go projects using Gherkin syntax. It's an alternative to [godog](https://github.com/cucumber/godog). There are a few things that differ and were my motivation to write my own library.

## How does GoBDD work?
You write your tests using the [gherkin](https://cucumber.io/docs/gherkin/) syntax. All tests (by default) should be placed in `features/` folder and have a prefix `.feature`.

The library reads all available documents. Loads all the step definitions. Next, tries to execute every scenario and steps into the scenario one by one. At the end, it produces the report of the execution.

You can have multiple gherkin documents executed within one test.

## What's the difference compared to godog?

Godog runs as an external binary. You have to download it and run as a separate step in your project. It generates a temporary test file, compiles it and runs in an external proces. This way of doings it has some disadvantages:

* working with the debugger is difficult. It runs in an subprocess, and it's not an easy task to attach the debugger to it.
* your code metrics are malformed. When you run tests with code coverage report, the code executed in godog doesn't count.
* some style checkers recognise tests as dead code.
* ...and a few more

That's why I decided to choose another approach and build GoBDD as an extension to built-in testing library. It uses a regular and well-known `testing.T`.

```go
func TestScenarios(t *testing.T) {
	suite := gobdd.NewSuite(t)
	// add your steps here
	suite.Run()
}
```

The second design decision I didn't like about godog is using the global state. Let's take a look at the example from their readme file.

```go
func thereAreGodogs(available int) error {
	Godogs = available
	return nil
}

func iEat(num int) error {
	if Godogs < num {
		return fmt.Errorf("you cannot eat %d godogs, there are %d available", num, Godogs)
	}
	Godogs -= num
	return nil
}

func thereShouldBeRemaining(remaining int) error {
	if Godogs != remaining {
		return fmt.Errorf("expected %d godogs to be remaining, but there is %d", remaining, Godogs)
	}
	return nil
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() { Godogs = 0 })
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {
		Godogs = 0 // clean the state before every scenario
	})

	ctx.Step(`^there are (\d+) godogs$`, thereAreGodogs)
	ctx.Step(`^I eat (\d+)$`, iEat)
	ctx.Step(`^there should be (\d+) remaining$`, thereShouldBeRemaining)
}
```

In those tests, the number of godogs is kept into a global variable. It means that it's not safe to run those tests concurrently and there are some indirect dependencies between those steps. To fix the problem I decided to add a context.

```go
func add(t gobdd.StepTest, ctx gobdd.Context, var1, var2 int) {
	res := var1 + var2
	ctx.Set("sumRes", res)
}

func check(t gobdd.StepTest, ctx gobdd.Context, sum int) {
	received, err := ctx.GetInt("sumRes")
    if err != nil {
        t.Error(err)

        return
    }

	if sum != received {
        t.Error(errors.New("the math does not work for you"))
	}
}

func TestScenarios(t *testing.T) {
	suite := NewSuite(t)
	suite.AddStep(`I add (\d+) and (\d+)`, add)
	suite.AddStep(`I the result should equal (\d+)`, check)
	suite.Run()
}

```

The context carriers the data from previously executed steps. It's concurrent safe. There's no need for creating a global variable. All you need lives in the `Context`.

GoBDD supports [parameter types](https://go-bdd.github.io/gobdd/parameter-types.html) and have a built-in set of the most popular once.

```go
    s := gobdd.NewSuite(t)
	s.AddParameterTypes(`{int}`, []string{`(\d)`})
	s.AddParameterTypes(`{float}`, []string{`([-+]?\d*\.?\d*)`})
	s.AddParameterTypes(`{word}`, []string{`([\d\w]+)`})
	s.AddParameterTypes(`{text}`, []string{`"([\d\w\-\s]+)"`, `'([\d\w\-\s]+)'`})
```

Thanks to this you can write

```go
suite.AddStep(`I add {int} and {int}`, add)
```
instead of 
```go
suite.AddStep(`I add (\d+) and (\d+)`, add)
```

This feature helps writing more readable tests. You can, of course, write your own parameter types using the `AddParameterTypes` function.

GoBDD has much more features and don't wait to try them. I'll be thankful for any kind of contribution - feeling a bug report, writing a suggestion or creating a PR to any [existing issue](https://github.com/go-bdd/gobdd/issues).

PS. If you like the project, tweet about it, share on FB or just tell your friends. I'll be thankful!