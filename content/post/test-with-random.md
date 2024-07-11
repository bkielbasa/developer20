---
title: "Testing code with unpredictable or random output"
publishdate: 2024-06-20
categories: 
    - Golang
    - Programming
tags:
  - golang
  - interview
  - testing
---

### Introduction

In this blog post, I want to share my approach to testing functions involving randomness in Go. Recently, I was asked how I would test a specific function that calculates possible directions for an object to move. Initially, I didn't come up with a good idea. Here, I'll discuss how I'd solve this problem in a real-world application with a detailed explanation.

### The Function in Question

The function calculates all possible directions that an object can move (up, down, left, right) without violating boundaries. It then randomly selects a valid direction and returns the new coordinates.

~~~go
type option struct {
    x, y int
}

const size = 10

func newPosition(x, y int) (int, int) {
    options := []option{
        {-1, 0},
        {1, 0},
        {0, -1},
        {0, 1},
    }

    possibilities := []option{}
    for _, opt := range options {
        newx := x + opt.x
        newy := y + opt.y

        if newx < 0 || newy < 0 {
            continue
        }

        if newx == size || newy == size {
            continue
        }

        possibilities = append(possibilities, opt)
    }

    n := rand.Intn(len(possibilities))
    r := possibilities[n]

    return x + r.x, y + r.y
}
~~~

### The TDD Problem

Test-Driven Development (TDD) is a powerful tool for improving code quality and increasing code coverage. While I don't use it daily, it's invaluable to have in my toolbox. The reason I don't use it every day is that TDD doesn't fit every situation. Testing the original function is a good example of where we need to seek better solutions.

### Dependency Inversion

A colleague suggested using dependency inversion to inject a function that would replace the `rand.Intn` call. This approach allows us to control the randomness during testing.

~~~go
func newPosition(x, y int, r func(int) int) (int, int) {
    options := []option{
        {-1, 0},
        {1, 0},
        {0, -1},
        {0, 1},
    }

    possibilities := []option{}
    for _, opt := range options {
        newx := x + opt.x
        newy := y + opt.y

        if newx < 0 || newy < 0 {
            continue
        }

        if newx == size || newy == size {
            continue
        }

        possibilities = append(possibilities, opt)
    }

    n := r(len(possibilities))
    r := possibilities[n]

    return x + r.x, y + r.y
}
~~~

This approach solves the problem, but there are two drawbacks:
1. It exposes a low-level detail about how the function works.
2. The testing code becomes more complicated and less readable due to mocking the `r(int) int` function.

### Statistical Testing

Instead of fighting the randomness of the behavior, we should accept it and focus on the function's behavior rather than its internal workings. My solution is inspired by how the Go team tests benchmarks and compares their results. They leave some room for unpredictability, and we should too.

~~~go
package main

import (
    "math/rand"
    "testing"
    "time"
)

func TestCalculate(t *testing.T) {
    rand.Seed(time.Now().UnixNano())

    // Define a grid for testing
    grid := make([][]int, size)
    for i := range grid {
        grid[i] = make([]int, size)
    }

    // Populate the grid with some test values
    grid[2][2] = 5
    grid[5][5] = 3

    iterations := 1000
    elementZeroCounts := make([]int, iterations)

    for i := 0; i < iterations; i++ {
        modified := calculate(grid)

        // Count elements with zero values
        zeroCount := 0
        for x := 0; x < size; x++ {
            for y := 0; y < size; y++ {
                if modified[x][y] == 0 {
                    zeroCount++
                }
            }
        }
        elementZeroCounts[i] = zeroCount
    }

    // Calculate average, min, max, and standard deviation
    var sum, sumSq float64
    min, max := elementZeroCounts[0], elementZeroCounts[0]

    for _, count := range elementZeroCounts {
        sum += float64(count)
        sumSq += float64(count * count)

        if count < min {
            min = count
        }
        if count > max {
            max = count
        }
    }

    mean := sum / float64(iterations)
    variance := sumSq/float64(iterations) - mean*mean
    stdDev := sqrt(variance)

    t.Logf("Mean: %v, Min: %v, Max: %v, StdDev: %v", mean, min, max, stdDev)

    // Check if the output is within an acceptable range
    if mean < 20 || mean > 40 {
        t.Errorf("Mean zero count out of expected range: %v", mean)
    }
}

// Simple square root function
func sqrt(x float64) float64 {
    z := x
    for i := 0; i < 1000; i++ {
        z -= (z*z - x) / (2 * z)
    }
    return z
}
~~~

By focusing on the statistical properties of the outputs over many runs, we ensure that the function behaves correctly without being overly deterministic about individual outputs.

### Conclusion

In this article, we've explored different approaches to testing functions involving randomness. While TDD is a powerful tool, it doesn't fit every situation. By using dependency inversion or focusing on statistical testing, we can ensure our functions work correctly without exposing low-level details or complicating our tests. We welcome your thoughts or alternative approaches in the comments section below!

