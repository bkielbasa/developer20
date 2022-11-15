---
title: "Writing TCP scanner in Go"
publishdate: 2019-10-23
categories: 
    - Golang
    - GoInPractice
    - Programming
tags:
  - golang
  - tcp
  - scanner
  - network
  - concurrency
resources:
    - name: header
    - src: featured.png

---

Go is perfect for network applications. Its awesome standard library helps a lot in writing such software. In this article, we’ll write a simple TCP scanner in Go. The whole programm will take less than 50 lines of code. Before we’ll go to practice - a little theory.

Of course, the TCP is [more complicated than I describe](http://www.medianet.kent.edu/techreports/TR2005-07-22-tcp-EFSM.pdf) but we need just basics. The TCP handshake is three-way. Firstly, the client sends the `syn` package which signals the beginning of a communication. If the client gets a timeout here it may mean that the port is behind a firewall.

![syn package](/images/diagram-01-sync.png)

Secondly, the server answers with `syn-ack` when the port is opened, otherwise it responses with `rst` package. In the end, the client has to send another packet called ack. From this point, the connection is established.

![syn package](/images/diagram-02-sync-ack.png)
![syn package](/images/diagram-03-ack.png)

The first step in writing the TCP scanner is to test a single port. We’ll use the `net.Dial` function which accepts two parameters: the protocol and the address to test (with the port number).

{{< highlight go >}}
package main

import (
	"fmt"
	"net"
)

func main() {
	_, err := net.Dial("tcp", "google.com:80")
	if err == nil {
		fmt.Println("Connection successful")
	} else {
		fmt.Println(err)
	}
}
{{< / highlight >}}

To not test every port one by one, we’ll add a simple loop that will simplify the whole process. Notice the Sprintf function which concretes the host and the port.

{{< highlight go >}}
package main

import (
	"fmt"
	"net"
)

func main() {
	for port := 80; port < 100; port++ {
		conn, err := net.Dial("tcp", fmt.Sprintf("google.com:%d", port))
		if err == nil {
			conn.Close()
			fmt.Println("Connection successful")
		} else {
			fmt.Println(err)
		}
	}
}
{{< / highlight >}}

The solution has one huge issue - it’s extremely slow. We can do two things to make things faster: run those checks concurrently and add a timeout to every connection.

Let’s focus on making in concurrent. The first step is to extract the scanning to a separate function. This step will make our code more clear.

{{< highlight go >}}
func isOpen(host string, port int) bool {
  conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
  if err == nil {
     _ = conn.Close()
     return true
  }

  return false
}
{{< / highlight >}}

The only new thing is the `WaitGroup`. You can read about it in more detail [here](https://gobyexample.com/waitgroups) or, if you want (let me know in the comments below) I can write an article about async programming in Go. But, back to the topic… In the main function, we span our goroutines and wait for the execution to finish.

{{< highlight go >}}
func main() {
  ports := []int{}

  wg := &sync.WaitGroup{}
  for port := 1; port < 100; port++ {
     wg.Add(1)
     go func() {
        opened := isOpen("google.com", port)
        if opened {
           ports = append(ports, port)
        }
        wg.Done()
     }()
  }

  wg.Wait()
  fmt.Printf("opened ports: %v\n", ports)
}
{{< / highlight >}}

Our code is faster but because of timeouts, we’re waiting a very long time to receive the error. We can assume that if we don’t get any response from the server for 200 ms we don’t want to wait longer.

{{< highlight go >}}
func isOpen(host string, port int, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err == nil {
		_ = conn.Close()
		return true
	}

	return false
}

func main() {
	ports := []int{}

	wg := &sync.WaitGroup{}
	timeout := time.Millisecond * 200
	for port := 1; port < 100; port++ {
		wg.Add(1)
		go func(p int) {
			opened := isOpen("google.com", p, timeout)
			if opened {
				ports = append(ports, p)
			}
			wg.Done()
		}(port)
	}

	wg.Wait()
	fmt.Printf("opened ports: %v\n", ports)
}

{{< / highlight >}}

At this point, we have a working simple port scanner. Unfortunately, it’s not very handy because to change the domain or port ranges we have to edit the code and recompile. Go has an awesome package called `flag`. 

The `flag` package helps in writing command-line applications. You can read more about it in [Go by Example](https://gobyexample.com/command-line-flags). What we want is configuring every magic string or number. We add parameters for the hostname, port range we want to test and the timeout on the connection.

{{< highlight go >}}
func main() {
	hostname := flag.String("hostname", "", "hostname to test")
	startPort := flag.Int("start-port", 80, "the port on which the scanning starts")
	endPort := flag.Int("end-port", 100, "the port from which the scanning ends")
	timeout := flag.Duration("timeout", time.Millisecond * 200, "timeout")
	flag.Parse()

	ports := []int{}

	wg := &sync.WaitGroup{}
	for port := *startPort; port <= *endPort; port++ {
		wg.Add(1)
		go func(p int) {
			opened := isOpen(*hostname, p, *timeout)
			if opened {
				ports = append(ports, p)
			}
			wg.Done()
		}(port)
	}

	wg.Wait()
	fmt.Printf("opened ports: %v\n", ports)
}
{{< / highlight >}}

If we want to show the usage, we have to put the -h parameter which will show us the usage. Simple and clear. The whole project took less than 50 lines of code. We used concurrency, the flag, and net packages.

There's one more thing. Our program has race condition. In only a few opened ports and so slow scanning it's not visible but there's the issue. To fix that we'll add [a mutex](https://gobyexample.com/mutexes).

{{< highlight go >}}
	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}
	for port := *startPort; port <= *endPort; port++ {
		wg.Add(1)
		go func(p int) {
			opened := isOpen(*hostname, p, *timeout)
			if opened {
				mutex.Lock()
				ports = append(ports, p)
				mutex.Unlock()
			}
			wg.Done()
		}(port)
	}
{{< / highlight >}}

If you like this kind of posts or have a question, let me know in the comments section below. The whole source code is available [on GitHub](https://github.com/bkielbasa/port-scanner).

