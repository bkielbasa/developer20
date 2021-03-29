---
title: "Writing custom linter in Go"
publishdate: 2021-03-31
resources:
    - name: header
    - src: featured.jpg
categories: 
    - Golang
    - Programming
tags:
  - golang
  - linter
---

Writing linters is simple. I was surprised how it's easy to write a Go linter. Today, we'll write a linter that will calculate the cyclomatic complexity of the Go code.

What is cyclomatic complexity?

> Cyclomatic complexity is a software metric used to indicate the complexity of a program. [ref](https://en.wikipedia.org/wiki/Cyclomatic_complexity)

The idea is simple - every time we find any control flow statements we increase the complexity by one. I know I oversimplified it a bit but I don't want to overwhelm you with unnecessary details. 

There are a few steps we should follow to write our custom linter. Firstly, we can create a test that will check if our linter works or not. Let's put him into `pkg/analyzer/analyzer_test.go` file.

```go
package analyzer

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata")
	analysistest.Run(t, testdata, NewAnalyzer(), "complexity")
}

```

The `analysistest.Run()` function is a helper that simplifies testing linters. What it does is running our linter on package `complexity` in `testate` folder. We use `NewAnalyzer()` function that will return instance of our analyzer. Let's add it to `pkg/analyzer/analyzer.go`.

```go
package analyzer

import "golang.org/x/tools/go/analysis"

//nolint:gochecknoglobals
var flagSet flag.FlagSet

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:  "cyclop",
		Doc:   "calculates cyclomatic complexity",
		Run:   run,
		Flags: flagSet,
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
  return nil, nil
}

```

We will use the flag set to input parameters to the linter. In the `NewAnalyzer()` we define the analyzer with the name `cyclop`, description, and defined `Run` function. It accepts `analysis.Pass` struct that's described [in the docs](https://pkg.go.dev/golang.org/x/tools@v0.1.0/go/analysis#Pass). We'll need only a few items from it right now.

To report an issue in the analyzed file we can use `pass.Reportf()` method. It accepts the position of the diagnostic and the message the user will see. The `pass.Files` is a slice of `*ast.File`, that is, a list of files within the package.

We'll use the last one to iterate every file and check them one by one. When we spot an issue - we'll report it. To do that, we have to iterate over those files using the loop below.

```go
for _, f := range pass.Files {
		ast.Inspect(f, func(node ast.Node) bool {
		  // your code goes here
		}
}

```

This is the heart of our linter. The `ast.Node` is an interface with only two methods that are needed while reporting any issues.

```go
type Node interface {
    Pos() token.Pos // position of first character belonging to the node
    End() token.Pos // position of first character immediately after the node
}

```

We're interested only in functions or methods so, we have to cast this type into `*ast.FuncDecl`.

```go
type FuncDecl struct {
    Doc  *CommentGroup // associated documentation; or nil
    Recv *FieldList    // receiver (methods); or nil (functions)
    Name *Ident        // function/method name
    Type *FuncType     // function signature: parameters, results, and position of "func" keyword
    Body *BlockStmt    // function body; or nil for external (non-Go) function
}

```

Our loop should look like the one below.

```go
for _, f := range pass.Files {
		ast.Inspect(f, func(node ast.Node) bool {
    f, ok := node.(*ast.FuncDecl)
			if !ok {
				return true
			}
	})
}

```

It's time to calculate the cyclomatic complexity. The new `complexity` function accepts the `*ast.File`. The algorithm increases the complexity every time it finds `if`, `for`, `select` or `case` statement as well as `||` or `&&` operators in `if` statement.

```go
func complexity(fn *ast.FuncDecl) int {
	v := complexityVisitor{Complexity:1}
	ast.Walk(&v, fn)
	return v.Complexity
}

type complexityVisitor struct {
	Complexity int
}

func (v *complexityVisitor) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.FuncDecl, *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause, *ast.CommClause:
		v.Complexity++
	case *ast.BinaryExpr:
		if n.Op == token.LAND || n.Op == token.LOR {
			v.Complexity++
		}
	}
	return v
}

```

There's `ast.Walk(&v, fn)` line that can be confusing. This function accepts [a visitor](https://golang.org/pkg/go/ast/#Visitor) that is called on every child node of the parent node we provide. We use a simple [type assertion](https://tour.golang.org/methods/15) to determinate the type of the specific node.

When we have the code, we can write test cases. To do that, let's create a new `complexity` package. Inside of `complexity.go` file put the code below.

```go
package complexity

import "testing"

func highComplexity() { // want "calculated cyclomatic complexity for function"
	i := 1
	if i > 2 {
		if i > 2 {
		}
		if i > 2 {
		}
		if i > 2 {
		}
		if i > 2 {
		}
	} else {
		if i > 2 {
		}
		if i > 2 {
		}
		if i > 2 {
		}
		if i > 2 {
		}
	}

	if i > 2 {
	}
}

func noComplexity() {}

```

We created two functions. The `noComplexity()` function has `1` complexity because it has only one execution path. The linter shouldn't report any problem. The `highComplexity` has a bunch of `if` statements that increase this metric. It **should** report an issue. Notice that we have a special comment that lets the test suite know about the diagnostic we expect to be reported. You can think about it as the `then` section in [given-when-then](https://martinfowler.com/bliki/GivenWhenThen.html) approach. Everything that's inside quotes is a regular expression that must match the diagnostic message. You can read more about it [in the official docs](https://pkg.go.dev/golang.org/x/tools@v0.1.0/go/analysis/analysistest#Run).

When we run our tests, we should see success! The next step is to make the linter executable. All we have to do is create a file in `cmd/cyclop/cyclop.go` file and put the content below.

```go
package main

import (
	"github.com/bkielbasa/cyclop/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.NewAnalyzer())
}

```

The code above runs our analyzer as a standalone tool. And that's it! We have a fully functional Go linter! The full (and a bit more complex) code is [available on Github](https://github.com/bkielbasa/cyclop).

There's one optional step you can make - add your linter to [golangci-lint](https://github.com/golangci/golangci-lint). Here's an example for the [linter we've built](https://github.com/golangci/golangci-lint/pull/1738).

I hope the Go community will get a lot of awesome linters that will save us hours or days.

