---
layout: post
title: Scientific method
mainPhoto: scientific-method.jpg
---

In 50′ and 60′ input data for programs from those years were written on paper tapes or punch cards. Writing code, compiling and testing loop took from a few hours to even few days. It was the beginning of programming we know it.

At this time Dijkstra started his discovery. He conceived the [algorithm called his name](https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm). He also noticed that programs become too complicated to be fully understood by one person. A single developer can miss something, forget or misunderstand. This situation could cause that program could work incorrectly. The program’s invalid behavior could be noticed after some time. In industries like financial or telecom it could lead to high costs, be hard to debug and fix.

To prevent such situations, Dijkstra proposed [mathematical discipline of proof](https://www.cs.utexas.edu/users/EWD/transcriptions/EWD03xx/EWD361.html). It means that a developer should build a theory which would solve a problem and try to prove its correctness mathematically. When the theory would be proven, he could write the code related to the mathematical statement. It would be the way the developers would be sure that their software works as expected.

Covering every part of the code with theory and then proving it was very time-consuming and complicated. It slowed down the development time which was already very slow. Developers found a different way to achieve that – a scientific method.

## Scientific method

A theory that the sun will rise every morning cannot be proven. On the other hand, the fact that the sun rises every day since we could take this information down is a strong enough argument to believe it will rise again.

Furthermore, it is impossible to prove that the sun will rise for the next 1000 years. On the other hand, the argument that it rises so long is good enough for us to be quite sure it will rise tomorrow. This is how the scientific method works.
Take a look at the example function below which implements [Fibonacci numbers algorithm](https://en.wikipedia.org/wiki/Fibonacci_number).

```python
function fibonacci(x)
    if x == 0 return 0
    if x == 1 return 1
    return fibonacci(x - 1) + fibonacci(x - 2)
```

A good way to test the above function is to add a set of tests which will cover most edge cases can be found and test its regular usage. In many situations, this set of tests is sufficient enough to be sure it works in any other case.

```
assert fibonacci(0) == 0
assert fibonacci(1) == 1
assert fibonacci(2) == 1
assert fibonacci(13) == 223
assert fibonacci(20) == 6765
```

To increase our trust in the implementation, more test cases can be provided using for example table tests with many well-known or previously calculated results of the function.

The scientific method can be used on many layers of the software. It is helpful in unit, integration, acceptance tests and so on. It is important to not cover only happy path but test edge cases too.

In the example above, the happy path is covered. But what if someone passes a negative value? How should it behave? Should it be ignored, raise an error or stop executing at all? It depends on the situation and the kind of the problem but they should be tested.

The disadvantage of the scientific method is that it does not give us in every case 100% certainty that the code we tested always works as expected. There is always a chance that a developer forgets or would not think about a specific edge case.

## Summary

Like almost everything, the scientific method has pros and cons but the fact is that it is a standard tool in the developers’ toolchain. The divide and conquer rule can cooperate with the scientific method to achieve the best balance between the time the development takes and assurance of working of stability of the application.