---
title: Cyclop
publishDate: 2021-12-13

---

[Cyclop](https://github.com/bkielbasa/cyclop) is a Go linter that calculates cyclomatic complexities of functions or packages in Go source code.

You can use it as a standalone application or using golangci-lint. I suggest to use the second one (it's more handy).

The default configuration looks like this:

```yaml
linters-settings:
  cyclop:
    # the maximal code complexity to report
    max-complexity: 10
    # the maximal average package complexity. If it's higher than 0.0 (float) the check is enabled (default 0.0)
    package-average: 0.0
    # should ignore tests (default false)
    skip-tests: false
```

Of course, don't forget to enable the linter.

```yaml
linters:
  enable:
    - cyclop
```

I found it very helpful when the function's lenght starts to grow. Setting some hard limits may be helpful to avoid creating 'monsters' in your code.