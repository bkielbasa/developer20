---
title: "Garnish - simple varnish implementation written in Go"
publishdate: 2020-02-03
categories: 
    - Golang
    - GoInPractice
    - Programming
tags:
  - golang
  - varnish
  - network
  - caching
resources:
    - name: header
    - src: featured.jpg
---

The varnish is a well-known HTTP accelerator. As the continuation of the [GoInPractice](https://developer20.com/categories/GoInPractice/) series, today I’ll show how you can build a simple (and naive) varnish implementation in Go. Some of the code is reused from [Writing a reverse proxy](https://developer20.com/writing-proxy-in-go/) so if you don’t understand something, I recommend taking a look at the blog post.

We’ll split our project into a few parts. The first one will be the caching mechanism. Its responsibility will be storing the data to cache and invalidate it after reaching the deadline. We’ll use [mutex](https://golang.org/pkg/sync/#Mutex)es for synchronization between goroutines.

Let’s start with tests. We’ll cover three of the most important scenarios: setting, and fetching data and two tests with timeouts. 

{{< highlight go >}}
func TestStoringAndRetrieving(t *testing.T) {
  c := newCache()
  data := []byte("data to store")
  c.store("key", data, 0)

  assert.Equal(t, data, c.get("key"))
}

func TestNotReachedTimeout(t *testing.T) {
  c := newCache()
  data := []byte("data to store")
  c.store("key", data, time.Millisecond*100)
  time.Sleep(time.Millisecond * 80)

  assert.Equal(t, data, c.get("key"))
}

func TestTimeout(t *testing.T) {
  c := newCache()
  data := []byte("data to store")
  c.store("key", data, time.Millisecond*100)
  time.Sleep(time.Millisecond * 100)

  assert.Equal(t, []byte(nil), c.get("key"))
}
{{< / highlight >}}

The cache keeps the data in a map. The most exciting part is the `store` function.

{{< highlight go >}}
func (c *cache) store(key string, rawData []byte, timeout time.Duration) {
  d := data{
     data: rawData,
  }
  if timeout != 0 {
     t := time.Now().Add(timeout)
     d.expires = &t
  }

  c.mutex.Lock()
  c.data[key] = d
  c.mutex.Unlock()

  time.AfterFunc(timeout, func() {
     c.clear(key)
  })
}
{{< / highlight >}}

At the beginning, we create the `data` struct which holds all the information about the cache: the data itself and the timeout (how long it stays in the memory). Then, we use mutexes to make sure we're not facing [race condition](https://en.wikipedia.org/wiki/Race_condition) problem and set the proper timeout.
 
But, it has some issues. For example,  when a new value is set to an existing key with higher duration then the previous key’s duration will invalidate the second value. Here's an example which demonstrates the issue.

{{< highlight go >}}
c.store(key, value, time.Millisecond*50)
c.store(key, value2, time.Millisecond*100) // this will be cleared after 50ms
{{< / highlight >}}

As long as we don’t have any other way of [invaliding the cache](https://www.varnish-software.com/wiki/content/tutorials/varnish/vcl_examples.html) we have nothing to worry about. You, the Reader, can add both functionalities as an excercise.

Our HTTP accelerator will support only one header - `Cache-Control`. There are a lot [headers related to caching](http://book.varnish-software.com/3.0/HTTP.html#cache-related-headers) but we won’t support all of them. Just in the name of simplicity.

Tests first. We’ll test it by running a HTTP server in the background which listen on port 8080 and return a response with `Cache-Control` header. We’ll send the request to the server twice and expect to receive two responses in `X-Cache` header: `MISS` and `HIT`.

{{< highlight go >}}
func TestGarnish_CacheRequest(t *testing.T) {
	stop := mockServer()
	defer stop()

	expectedXCacheHeaders := []string{garnish.XcacheMiss, garnish.XcacheHit}
	g := garnish.New(url.URL{Scheme: "http", Host: "localhost:8088"})

	for _, expectedHeader := range expectedXCacheHeaders {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8088", nil)
		w := httptest.NewRecorder()
		g.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		xcache := w.Header().Get("X-Cache")
		assert.Equal(t, expectedHeader, xcache)
	}
}
{{< / highlight >}}

We check if the HTTP method supports caching. Only GET requests should be cached. Then, we use [the reverse proxy](https://developer20.com/writing-proxy-in-go/) to pass the request further. When the response is returned we take a look at the `Cache-Control` header and based on this information we make the decision: to cache or not to cache. In the end, we add the `X-Cache` header which informs the client if the response was cached or not.

{{< highlight go >}}
func (g *garnish) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// only GET requests should be cached
	if r.Method != http.MethodGet {
		rw.Header().Set(Xcache, XcacheMiss)
		g.proxy.ServeHTTP(rw, r)
		return
	}

	u := r.URL.String()
	cached := g.c.get(u)
	if cached != nil {
		rw.Header().Set(Xcache, XcacheHit)
		_, _ = rw.Write(cached)
		return
	}

	proxyRW := &responseWriter{
		proxied: rw,
	}
	proxyRW.Header().Set(Xcache, XcacheMiss)
	g.proxy.ServeHTTP(proxyRW, r)

	cc := rw.Header().Get(cacheControl)
	toCache, duration := parseCacheControl(cc)
	if toCache {
		g.c.store(u, proxyRW.body, duration)
	}
}
{{< / highlight >}}

The whole source code of the project is availabe [on Github](https://github.com/bkielbasa/garnish) so you can experiment with it by your own. Adding new features should be straightforward.  We needed only the standard library. If you have any questions, let me know in the comments section below.
