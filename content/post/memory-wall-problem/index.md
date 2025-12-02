---
title: "Memory-wall problem"
publishdate: 2023-06-01
categories: 
    - Programming
tags:
  - golang
  - memory
  - performance
---

The memory wall problem refers to a phenomenon that occurs in computer architecture when the processor's speed outpaces the rate at which data can be transferred to and from the memory system. As a result, the processor must wait for the data to be fetched from memory, which slows down its performance and limits its speed.

The memory wall problem has become increasingly significant as processors have become faster and more powerful while memory speeds have not kept pace with these advancements. This means that even though the processor can execute instructions quickly, it spends a significant amount of time waiting for data to be transferred to and from memory.

To mitigate the problem, CPUs have special L-caches that are small but very fast. They are static random-access memory (SRAM) that are much faster than Main Memory but are much more expensive. Their purpose is to hold small, frequently used data for quick access. Firstly, the CPU tries to read the data from the L1 cache, if it can't find it, checks the L2 cache, and so on. If it doesn't find it, it reads from the main memory. The Main Memory is our RAM that is the biggest but the slowest at the same time. In the fixture below, you can find sample PC storages with sizes and access time.

{{< figure src="https://i.imgur.com/suD65Rs.png" title="Image source: https://cs.brown.edu/courses/csci1310/2020/assign/labs/lab4.html" >}}

As you can see, the smaller the memory size is, the faster the read access. In comparison, reading the same data from L1-cache is about 60 times faster than reading the same data from the main memory. In every PC those values will differ, depending on the hardware. You can check L-cache sizes using one of the following commands.

For linux users:

```sh
lscpu

# Sample output:

Architecture:          x86_64
CPU op-mode(s):        32-bit, 64-bit
Byte Order:            Little Endian
CPU(s):                4
On-line CPU(s) list:   0-3
Thread(s) per core:    1
Core(s) per socket:    4
Socket(s):             1
NUMA node(s):          1
Vendor ID:             GenuineIntel
CPU family:            6
Model:                 42
Stepping:              7
CPU MHz:               3401.000
BogoMIPS:              6784.57
Virtualization:        VT-x
L1d cache:             32K
L1i cache:             32K
L2 cache:              256K
L3 cache:              8192K
NUMA node0 CPU(s):     0-3
```

For macos users:

```sh
sysctl -a | grep hw.l

# Sample output:

hw.logicalcpu: 10
hw.logicalcpu_max: 10
hw.l1icachesize: 131072
hw.l1dcachesize: 65536
hw.l2cachesize: 4194304
```

{{< question title="What do L1i and L1d mean?" msg="There are two types of L1 caches. L1i stands for instruction cache (used to speed up executable instruction fetch) and L1d for data cache (used to speed up data fetch and store)." >}}

The CPU doesn't read a value one by one but does cache lines. Of course, depending on the hardware the cache line size may differ as well. A single cache line reads the requested value plus next values in the memory. Reading a single integer or, for example, 8 of them at once takes exactly the same amount of time so it does it with a hope to read what we probably will need next. It's a good way to limit cache-misses.

The most useful data type are slices and arrays. In a loop that iterates over an array, the program is likely to access nearby elements of the array repeatedly. Storing these elements in the cache can significantly improve the program's performance.

We can check how the memory is allocated using a loop and a `%p` directive. Below, you can find examples for booleans and integers.


```go
func main() {
	s := []byte{1, 0, 1, 0}

	for i := 0; i < len(s); i++ {
		fmt.Printf("%v\t%p\t%d bytes\n", s[i], &s[i], unsafe.Sizeof(s[i]))
	}
}

/*
prints:

1	0xc0000b2000	1 bytes
0	0xc0000b2001	1 bytes
1	0xc0000b2002	1 bytes
0	0xc0000b2003	1 bytes
*/
```

```go
func main() {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8}

	for i := 0; i < len(s); i++ {
		fmt.Printf("%v\t%p\t%d bytes\n", s[i], &s[i], unsafe.Sizeof(s[i]))
	}
}

/*
prints:

1	0xc0000b2000	8 bytes
2	0xc0000b2008	8 bytes
3	0xc0000b2010	8 bytes
4	0xc0000b2018	8 bytes
5	0xc0000b2020	8 bytes
6	0xc0000b2028	8 bytes
7	0xc0000b2030	8 bytes
8	0xc0000b2038	8 bytes
*/
```
{{< emailform>}}

Every datum is stored one by one, without any gaps. Reading a value in `s[0]` will end up reading 4-8 next values at once which makes things much faster.

In comparison, lists have pointers that may point to memory that's not necessary in the cache line. It will end up with more cache-miss rate and slower programs.

```go
package main

import (
	"fmt"
	"time"
)

type node struct {
	value int
	next  *node
}

func generate() *node {
	root := &node{
		value: 200,
	}

	prev := root

	for i := 0; i < 10; i++ {
		n := node{
			value: 300 + i,
		}
		prev.next = &n
		prev = &n
		time.Sleep(time.Millisecond * 100)
	}

	return root
}

func main() {
	var n1, n2 *node
	ch := make(chan struct{})

	go func() {
		n1 = generate()
		ch <- struct{}{}
	}()

	go func() {
		n2 = generate()
		ch <- struct{}{}
	}()

	<-ch
	<-ch

	root := &node{
		value: 100,
		next:  n1,
	}

	fmt.Println("root")
	printNodes(root)

	fmt.Println("n2")
	printNodes(n2)

}

func printNodes(n *node) {
	for n != nil {
		fmt.Printf("val: %d, memory: %p\n", n.value, n)
		n = n.next
	}
}

/*
prints:

root
val: 100, memory: 0x140000100d0
val: 200, memory: 0x14000100000
val: 300, memory: 0x14000100010
val: 301, memory: 0x14000100020
val: 302, memory: 0x14000100030
val: 303, memory: 0x14000100040
val: 304, memory: 0x14000100050
val: 305, memory: 0x14000100060
val: 306, memory: 0x14000100070
val: 307, memory: 0x14000100080
val: 308, memory: 0x14000100090
val: 309, memory: 0x140001000a0
n2
val: 200, memory: 0x14000096020
val: 300, memory: 0x14000096030
val: 301, memory: 0x14000010040
val: 302, memory: 0x14000010050
val: 303, memory: 0x14000010060
val: 304, memory: 0x14000010070
val: 305, memory: 0x14000010080
val: 306, memory: 0x14000010090
val: 307, memory: 0x140000100a0
val: 308, memory: 0x140000100b0
val: 309, memory: 0x140000100c0
*/
```

As you can see, especially at the beginning, those values aren't next to each other. In a real-world application those values will be much more random. This situation happens because creating a new node in the list is a separate allocation. When creating a slice, it's just one.

To illustrate how it will impact the performance, let’s discuss the following program:

```go
package caching

import "fmt"

// Create a square matrix of 16,777,216 bytes.
const (
	rows = 4 * 1024
	cols = 4 * 1024
)

// matrix represents a matrix with a large number of
// columns per row.
var matrix [rows][cols]byte

// data represents a data node for our linked list.
type data struct {
	v byte
	p *data
}

// list points to the head of the list.
var list *data

func init() {
	var last *data

	// Create a link list with the same number of elements.
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {

			// Create a new node and link it in.
			var d data
			if list == nil {
				list = &d
			}
			if last != nil {
				last.p = &d
			}
			last = &d

			// Add a value to all even elements.
			if row%2 == 0 {
				matrix[row][col] = 0xFF
				d.v = 0xFF
			}
		}
	}

	// Count the number of elements in the link list.
	var ctr int
	d := list
	for d != nil {
		ctr++
		d = d.p
	}

	fmt.Println("Elements in the link list", ctr)
	fmt.Println("Elements in the matrix", rows*cols)
}

// LinkedListTraverse traverses the linked list linearly.
func LinkedListTraverse() int {
	var ctr int

	d := list
	for d != nil {
		if d.v == 0xFF {
			ctr++
		}

		d = d.p
	}

	return ctr
}

// ColumnTraverse traverses the matrix linearly down each column.
func ColumnTraverse() int {
	var ctr int

	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			if matrix[row][col] == 0xFF {
				ctr++
			}
		}
	}

	return ctr
}

// RowTraverse traverses the matrix linearly down each row.
func RowTraverse() int {
	var ctr int

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if matrix[row][col] == 0xFF {
				ctr++
			}
		}
	}

	return ctr
}
```

You can find a matrix of 16,777,216 bytes in it. Our goal is to iterate over it using 3 ways: rows first, columns first and using a list. In the `init()` function you can find allocating the list so its creation won’t impact our benchmarks. Let’s write some benchmarks for those functions.

```go
package caching

import "testing"

var fa int

func BenchmarkLinkListTraverse(b *testing.B) {
	var a int

	for i := 0; i < b.N; i++ {
		a = LinkedListTraverse()
	}

	fa = a
}

func BenchmarkColumnTraverse(b *testing.B) {
	var a int

	for i := 0; i < b.N; i++ {
		a = ColumnTraverse()
	}

	fa = a
}

func BenchmarkRowTraverse(b *testing.B) {
	var a int

	for i := 0; i < b.N; i++ {
		a = RowTraverse()
	}

	fa = a
}
```

We use the global variable to avoid some compiler optimizations that can affect our tests. When we run our benchmarks we’ll notice that the difference between those results are significant.


```sh
$ go test -run none -bench . -benchtime 3s

Elements in the link list 16777216
Elements in the matrix 16777216
goos: darwin
goarch: arm64
pkg: caching
BenchmarkLinkListTraverse-10                 207          17241430 ns/op
BenchmarkColumnTraverse-10                    74          46370815 ns/op
BenchmarkRowTraverse-10                      418           8584940 ns/op
```

The fastest function is the one that traverses the matrix's rows first. It happens because it has the smallest cache-miss ratio. The linked list implementation is about 3 times slower but it’s still faster than the slowest row-first travers that’s about 5.6 times slower. That’s a notable difference.

## Summary

The memory-wall problem presents a significant challenge in maximising application performance. As always, measuring the performance is the key and you should never guess if a code change will impact the performance in any way. At some point, it's very useful to understand how computers work and how we can help them to work faster. Memory-wall problem is only one of challenges we have in this field.

Of course, I only scrached the surface of profiling the code so if you want more content about it - let me know in the comments section below.
