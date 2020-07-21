---
title: "Add header to every request in Go"
publishdate: 2020-07-21
resources:
    - name: header
    - src: featured.jpg
---

Making changes to all HTTP requests can be handy. You may want to add an API key or some information about the sender like app version etc. No matter why you want to do that you have a few options to achieve the goal.

The first approach is building a factory method that will add the required header.


{{< highlight golang >}}
func newRequest(endpoint string) *http.Request {
    req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s", endpoint), nil)
    req.Header.Add("x-api-key", "my-secret-token")
    return req
}
{{< / highlight >}}

It’s very simple and clear but it requires you to create a new method per HTTP method or calling another function directly. You’ll have to do it every time you create a new request.

{{< highlight golang >}}
func addHeaders(r *http.Request) {
	req.Header.Add("x-api-key", "my-secret-token")
}

// somewhere in your code

req, err := http.NewRequest(/* ... */)
if err != nil {
	return nil
}

addHeaders(req)

// do your job here

{{< / highlight >}}

This approach is simple but requires repeating the same code over and over, across all the modules that can be error-prone. You can achieve a similar result much easier - using [RoundTripper](https://godoc.org/net/http#RoundTripper).

{{< highlight golang >}}
type transport struct {
	underlyingTransport http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("x-api-key", "my-secret-token")
	return t.underlyingTransport.RoundTrip(req)
}
{{< / highlight >}}

and while creating the HTTP client

{{< highlight golang >}}
c := http.Client{Transport: &transport{ underlyingTransport: http.DefaultTransport } }
{{< / highlight >}}

This looks elegant, clear and simple. However, there's a comment in the docs

>  RoundTrip should not modify the request, except for consuming and closing the Request's Body.

It means we **shouldn't** mutate the request the way we did it. The round trip is often used for rate limiting or caching responses and it's more valid usage. Is there any other way of doing that? Of course - use the power of Go interfaces.

What you should do is to write a wrapper/middleware for the `http.Client` with all of the methods you use and add your headers there.

{{< highlight golang >}}
type httpClient struct {
	c        http.Client
	apiToken string
}

func (c *httpClient) Get(url string) (resp *Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err

	}

	return c.Do(req)

}

func (c *Client) Post(url, contentType string, body io.Reader) (resp *Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err

	}

	req.Header.Set("Content-Type", contentType)

	return c.Do(req)

}

func (c *httpClient) Do(req *Request) (*Response, error) {
	req.Header.Add("x-api-key", c.apiToken)
	return c.c.Do(req)

}
{{< / highlight >}}

As you can see, it requires more code but on the other hand, it gives you more flexibility about where and what we change in the request.

## Summary
There are a few ways you can achieve the same or similar results. Which one should you choose? It depends on your use case. Always try to find the simples and a more readable way. If the factory method looks fine for you, don’t overcomplicate your code. When your factory method is becoming more and more complicated, consider refactoring it into smaller methods or just use the wrapper.

