
---
title: "Pointer and value semantics in Go"
categories: [Golang, Programming]
publishDate: 2020-06-29
resources:
    - name: header
    - src: featured.jpg

tags:
 - golang
---

In Go, we can refer to variables using value or pointers. Sometimes, it's hard to answer which approach is more suitable.

At the first place, you should learn about general rules. Value semantic should be used every time when copying the value make sense in the logic of your code. For example, every value object should be passed by value. If you have a struct `Money` then it's possible (and also make sense) to have, at the same time, multiple 10$ in your code. There's no such requirement that only one person can have the same amount of money.

On the other hand, pointer semantic should be used in the opposite scenarios. When you have `Order` with a specific `ID`, you don't want to allow copying the `Order` object because the caller may not see changes in the copy.


{{< highlight golang >}}
func serviceCall(o order) {
    // do some logic here
    o.status = "processing"
    anotherServiceCall(o)
        fmt.Print(o.status) // prints "processing"
}

func anotherServiceCall(o order) {
    o.status = "closed"
}
{{< / highlight >}}

This can be even more unpredictable (at least at the beginning) for reference types inside structs like maps, pointers and slices because in this case, the caller **will** see the changes. Some changes will be visible outside of the `anotherServiceCall` function and some won't. It's easy to forget about it.

In other words, there can be multiple copies of the same book (value semantic) but only one person with the same ID number as you :) If you use value semantics, you copy the value. When using pointers - share it with others.

Another thing worth remembering is that if you decide to one semantic within one struct - stick to it. It should be clear to you and other devs in your team which approach is in use right now. Personally, I introduced a few bugs by forgetting to change from the value to pointer semantic in one of the methods. Such mistakes can be annoying.

## What about built-in types?

n Go, built-in types like integers, floats and so on should be passed as values. There are exceptions when you want explicitly share with others. An excellent example of this exception is `sql.Scan()` function that expects pointers. If you would pass the variable by value, it would (it won't because it checks if it's a pointer at the very beginning) override the copy of your value. The original variable would be unchanged.

When it comes to structs, they should be passed as values as long as it fits the general rule. If it makes sense to create multiple copies of the variable - use value semantics. If no, use pointers. But... there are exceptions as well.

Use pointers if the copy of the struct is expensive. One of the reason can be a very large struct or creating a lot of copies. In such cases, allocations and the GC overhead can have a significant impact on performance and memory usage.

What's more, it's OK to mix value and pointers when creating functions for deserialising. Let's consider the following example.


{{< highlight golang >}}
type Money struct {
  quantity int
	curr string
}

func (m Money) Add(toAdd Money) Money {
	if m.curr != toAdd.curr {
		panic("cannot add money with two different currencies")
  }

  return Money{
		curr: m.curr,
		Quantity: m.quantity + toAdd.quantity
	}
}

func (m Money) Currency() string {
	return m.curr
}

func (m Money) Quantity() int {
	return m.quantity
}

func (m Money) Valid() bool {
	return m.Quantity >= 0 && m.Currency != ""
}

func (m *Money) UnmarshalJSON(data []byte) error {
    // the implementation
}
{{< / highlight >}}

Functions like `Add()`, `Currency()`or `Valid()` accept the copy of the `Money` value object. Those methods don't have to have the write access to the original struct. In such case, it makes more sense to create a new struct instead of modifying the existing one.

## Does Go have reference types?

In Go, there's no "real" pass-by-reference. It is not possible to create two variables that share the same memory. For example, in c++ it's possible and quite easy.

{{< highlight cpp >}}
#include <iostream>

int main() {
    int p = 123;
    int &i = p; // the reference
    
    std::cout << &i << std::endl << &p << std::endl;
    return 0;
}
{{< / highlight >}}

In my previous article about [slices in Go](https://developer20.com/what-you-should-know-about-go-slices/) I described that every slice, in fact, is a struct with three fields: length, capacity and pointer to a backing array. When you pass your slices or arrays, those values are copied. We can argue if Go has reference types or not but it's worth remembering how it works.

## Summary

As you can see, deciding whatever to use or not pointers should be quite straightforward and reasonable. The copy/not-copy decision is strongly inspired by Domain-Driven Desing (DDD). We can argue if Go has [reference types or not](https://github.com/go101/go101/wiki/About-the-terminology-%22reference-type%22-in-Go), both semantics have pros and cons and should be used where appropriate.
