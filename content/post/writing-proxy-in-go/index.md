---
title: "Writing a reverse proxy in Go"
publishdate: 2019-11-25
featuredImage: /images/proxy.jpg
resources:
    - name: header
    - src: featured.jpg
categories: 
    - Golang
    - GoInPractice
    - Programming
tags:
  - golang
  - tcp
  - scanner
  - network
---

Some time ago, I found a video called [Building a DIY proxy with the net package](https://www.youtube.com/watch?v=J4J-A9tcjcA). I recommend watching it. Filippo Valsorda builds a simple proxy using low-level packages. It’s fun to watch it but I think it’s a bit complicated. In Go, it has to be an easier way so I decided to continue writing series [https://developer20.com/categories/GoInPractice/](Go In Practice) by writing a simple but yet powerful reverse proxy as fast as it’s possible.

The first step will be to create a proxy for a single host. The core of our code will be [https://golang.org/pkg/net/http/httputil/#ReverseProxy](ReversProxy) which does all the work for us. This is the magic of the rich standard library. The `RevewseProxy` is a struct for writing reverse proxies :) The only thing we have to do is to configure the director. The director modifies original requests which will be sent to proxied service.

{{< highlight go >}}
package main

import (
    "flag"
    "fmt"
    "net/http"
    "net/http/httputil"
)

func main() {
      url, err := url.Parse("http://localhost:8080")
      if err != nil {
          panic(err)
      }

    port := flag.Int("p", 80, "port")
    flag.Parse()

    director := func(req *http.Request) {
        req.URL.Scheme = url.Scheme
        req.URL.Host = url.Host
    }

    reverseProxy := &httputil.ReverseProxy{Director: director}
      handler := handler{proxy: reverseProxy}
    http.Handle("/", handler)

    http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

type handler struct {
    proxy *httputil.ReverseProxy
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.proxy.ServeHTTP(w, r)
}
{{< / highlight >}}

And that’s all! We have a fully functional proxy! Let’s check if it works. For the test, I wrote a simple server that returns the port it’s listening on.

{{< highlight go >}}
package main

import (
    "flag"
    "fmt"
    "net/http"
)

func main() {
    var port = flag.Int("p", 8080, "port")
    flag.Parse()
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(fmt.Sprintf("hello on port %d", *port)))
    })
    err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
    fmt.Print(err)
}
{{< / highlight >}}

Now, we can run the service which listens on port 8080 and the proxy which listens on port 80.

![Working TCP proxy](/images/proxy01.png)
Our proxy works on HTTP only. To tweak it to support HTTPS we have to make a small change. The thing we need to do is to detect (in a very naive way) if the proxy is running using SSL or not. We’ll detect it based on the port it’s running on.

{{< highlight go>}}
if *port == 443 {
    http.ListenAndServeTLS(fmt.Sprintf(":%d", *port), "localhost.pem", "localhost-key.pem", handler)
} else {
    http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
{{< / highlight >}}

To make it working, we need one more thing: a valid certyficate. You can generate it using [https://github.com/FiloSottile/mkcert](mkcert).

```
mkcert localhost
```

And... that’s all! We have a proxy that works on both HTTP/HTTPS.
![Working TCP proxy](/images/proxy02.png)

As you can see, writing the new tool was extremely simple and all we need is the standard library. We didn't have to write complicated code so we'll can focus on our real goals. As usual, the source code is available on Github: https://github.com/bkielbasa/go-proxy. This project will be a starting point for other tools so keep in touch and don’t miss a post! If you have any questions, let me know in the comments section below.
