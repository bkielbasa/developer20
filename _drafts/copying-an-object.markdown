---
layout: post
title:  Copying an object
mainPhoto: History-of-WWW.jpeg
---

In one of my previous articles, I wrote about [copying in golang]({% post_url 2019-01-25-be-aware-of-coping-in-go %}). After publishing this blog post, I received feedback that a similar error is very common in Python when someone starts learning the language. I'm sure that not only in those two languages it can be reproduced very easily. In this article, I'll show you the same kind of errors in other languages.

## Python example

We'll start with Python because it's an interesting example.

```python
class SomeClass:
    name = ""
    arrs = []


if __name__ == '__main__':
    c1 = SomeClass()
    c1.name = "my first array"
    c1.arrs.append(123)
    c1.arrs.append(321)
    c1.arrs.append(231)

    c2 = c1
    c2.name = "my second array"

    c2.arrs.append(777)

    print(c1.name)
    print(c1.arrs)
```

When you run the program you'll notice that both `name` and `arrs` parameters are overridden. This situation happens because in Python the object is not copied but a new reference to it is added. To change this behavior, we have to use `import` library.

```python
from copy import copy


class SomeClass:
    name = ""
    arrs = []


if __name__ == '__main__':
    c1 = SomeClass()
    c1.name = "my first array"
    c1.arrs.append(123)
    c1.arrs.append(321)
    c1.arrs.append(231)

    c2 = copy(c1)
    c2.name = "my second array"

    c2.arrs.append(777)

    print(c1.name)
    print(c1.arrs)
```

The name did not change but the `arrs` attribute did. It won't change if you use `deepcopy()` function. Unfortunately, it's not possible to prevent from bugs like that in Python. Very similar to the [Go example]({% post_url 2019-01-25-be-aware-of-coping-in-go %}).

## Java example
Let's take a look at some Java code. We have a very similar example: a class with an array attribute. We add a few elements to the array and then copy the whole object. At the end, we add a new element to the new object and print the FIRST object to see his content.

```java
import java.util.ArrayList;

class SomeClass {
    private ArrayList<Integer> myList;

    SomeClass() {
        myList = new ArrayList<>();
    }

    void add(int i) {
        this.myList.add(i);
    }

    void print() {
        System.out.println("********");
        for (int i = 0; i < myList.size(); i++) {
            System.out.println(myList.get(i));
        }
    }
}

public class Coping {
    public static void main(String[] args) {
        SomeClass l1 = new SomeClass();
        l1.add(123);
        l1.add(321);
        l1.add(231);

        l1.print();

        SomeClass l2 = l1;

        l2.add(777);
        l1.print();
    }
}
```

As before, the first object which should have only 3 elements but has 4 of them.  

## C# code

```
using System.Collections.Generic;

class SomeClass
{
	public List<int> arrs { get; set; }
	
	public SomeClass() {
		this.arrs = new List<int>();
	}
	
	public void Print() {
		foreach (int i in arrs)
		{
			System.Console.Write("{0}\n", i);
		}
	}
}
					
public class Program
{
	public static void Main()
	{
		SomeClass c1 = new SomeClass();
		c1.arrs.Add(123);
		c1.arrs.Add(321);
		c1.arrs.Add(231);
		
		SomeClass c2 = c1;
		
		c2.arrs.Add(777);
		
		c1.Print();
	}
}
```

