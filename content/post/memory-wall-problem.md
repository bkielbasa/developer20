---
title: "Memory wall problem and how Go helps to solve it"
publishdate: 2023-05-05
categories: 
    - Programming
tags:
  - golang
  - memory
---

The memory wall problem refers to a phenomenon that occurs in computer architecture when the processor's speed outpaces the rate at which data can be transferred to and from the memory system. As a result, the processor must wait for the data to be fetched from memory, which slows down its performance and limits its speed.

The memory wall problem has become increasingly significant as processors have become faster and more powerful while memory speeds have not kept pace with these advancements. This means that even though the processor can execute instructions quickly, it spends a significant amount of time waiting for data to be transferred to and from memory.

To mitigate the problem CPUs have special L-cache that are small but very fast. They are static random-access memory (SRAM) that are much faster than Main Memory but are much more expensive. Theier purpose is to hold small, frequently used data for quick access. Firstly, the CPU tries to read the data from L1 cache, if couldn't find checks the L2 cache, and so on. If it didn't find it, it reads from the main memory. The Main Memory is our RAM that is the biggest but the slowest at the same time. In the fixture below, you can find a sample PC storages with sizes and access time.

{{< figure src="https://i.imgur.com/suD65Rs.png" title="Image source: https://cs.brown.edu/courses/csci1310/2020/assign/labs/lab4.html" >}}

As you can see, the smaller the memory size is, the fastest the read access. In comparision, reading the same data from L1-cache is about 60 times faster than reading the same data from the main memory. In every PC those values will differ, depending on the hardware. You can check L-cache sizes using one of following commands.

For linux users

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

For macos users

```sh
sysctl -a | grep hw.l

# Sample output:

hw.logicalcpu: 10
hw.logicalcpu_max: 10
hw.l1icachesize: 131072
hw.l1dcachesize: 65536
hw.l2cachesize: 4194304
```

The CPU doesn't read a value one by one but does it it cache lines. Of course, depending on the hardware the cache line size may differ as well. A single cache line reads the requested value plus next values in the memory. Reading a single integers or, for example, 8 of them at once takes exactly the same amount of time so it does it with a hope to read what we probably will need next. It's a good way to limit cache-misses.

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

Every datum is stored one by one, without any gaps. Reading a value in `s[0]` will end up reading 4-8 next values at once what makes things much faster.
