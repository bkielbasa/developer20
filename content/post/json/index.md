---
title: "JSON in Go"
publishdate: 2022-10-17
tags:
  - json
  - reflect
toc: true
---

In this article, I'll tell you everything that you need to start using JSON in Go Fluent. We'll start with some basic usage, we'll talk about different ways of working in JSON and how to customize it. In the end, you'll find a FAQ where I put the most common questions about JSON in Go.

{{< table_of_contents >}}

## Basic usage

### Marshaling
Go has a built-in mechanism for marshaling (encoding) and unmarshaling (decoding) JSON files. The functionality is paced into `encoding/json` package. The basic usage is as follows:

```go
data, err := json.Marshal(yourVar)
```

The `Marshal()` method accepts any type that we want to marshal and returns a `[]byte` and an `error`. The data is ready JSON we can use if the `err` is a `nil`.

```go
data, err := json.Marshal(yourVar)
if err != nil {
  return err
}

// we can use `data` without worries
```

The error is returned only when the type we pass to the `Marshal()` method cannot be correctly encoded. Types that return the `*json.UnsupportedTypeError` are:

* channels
* complex [^complex]
* function values

```go
ch := make(chan struct{})
_, err := json.Marshal(ch) // returns error

compl := complex(10, 11)
_, err = json.Marshal(compl) // returns error

fn := func() {}
_, err = json.Marshal(fn) // returns error
```

It means that if a struct or a map contains those values, we'll get the error.

{{< info title=" Which fields in a struct are visible to the `json` package?" msg="It's a common mistake at the very beginning to forget that only **public struct fields** are used by the `json` package. It means that if the struct has a private field it won't be both marshaled and unmarshalled" >}}


### Unmarshalling

When we get a `[]byte` with our JSON we can easily decode it into our type thanks to the `json.Unmarhal()` method.

```go
myVal := MyVal{}
byte := `{"some":"json"}`
err := json.Unmarhal(byte, &myVal)
```

The error will be returned in the following cases:
* the data isn't a valid JSON
* we didn't provide the pointer to our local variable
* we provide a `nil` as the second parameter

{{< info title="Passing the pointer" msg="It's a common mistake to forget to add the pointer. Your IDE or linters may help you with catching such bugs." >}}

Go unmarshals the data into struct fields using either the struct field name or its tag. If it won't find it, it will try the case-insensitive match. 

## Struct tags
We can use struct tags to manipulate the way how fields are named in your JSON out or change mapping them in unmarshalling. Let me explain it in more detail.

Let's say we have a struct with two fields as shown below. When we encode the struct into JSON both fields will be capitalized. Very often, it's now what we want.

```go
type User struct {
  ID string
  Username string
}

// the output may look like this:
{"ID":"some-id","Username":"admin"}
```

To change the behavior we can use struct tags. After the field type, we add text. The first word it's the field tag name `json`, after it, we put `:` and in double quotes value of the tag. You can see an example below.

```go
type User struct {
  ID string `json:"id"`
  Username string `json:"user"`
}

u := User{ID: "some-id", Username: "admin"}

// the output may look like this:
{"id":"some-id","user":"admin"}
```

In the example, we renamed both fields. The name can be anything that's a valid JSON key. The standard library gives us one additional option: `omitempty`. We add it to fields that should be skipped if its value is `false`, `0`, a `nil` pointer, a `nil` interface value, and any empty array, slice, map, or string. We specify options after the JSON key and separate them with a comma (`,`).

```go
type User struct {
  ID string `json:"id"`
  Username string `json:"user"`
  Age string `age,omitempty`
}
```

If we don't want to change the default field name, we can skip it. We have to remember that in that case, the comma should be there anyway.

```go
type User struct {
  ID string `json:"id"`
  Username string `json:"user"`
  Age string `json:",omitempty"` // don't forget about the comma
}
```

If we want to keep the field public but tell the marshaller/unmarshaller to ignore it, we have to put a `-` to the tag value.

```go
type User struct {
  ID string `json:"id"`
  Username string `json:"user"`
  Age string `json:"-"`
}

u := User{ID: "some-id", Username: "admin", Age: 18}

// the output looks like this (notice missing age):
{"id":"some-id","user":"admin"}
```

{{< warning title="The struct tags are evaluated in runtime" msg=" If the runtime will have problems with parsing them (a parse error) the compiler won't complain about it. It's error-prone so it's important to remember about it" >}}

## Encoder/decoder
There are also `json.Decoder` and `json.Encoder` in the `json` package. They work similar to `json.Marshal()` and `json.Unmarshal()` methods. The biggest difference is that the first pair works on `io.Reader` and `io.Writer`. The second pair (marshal/unmarshal) work on a slice of bytes. It means it's more handy to use `json.Decoder`/`json.Encoder` if we don't have the data yet. I prepared two simple tables that should help us understand which option we should use.

When we decode data:

|  | `[]byte` | `io.Reader` |
|:--|:--|:--|
| `json.Unmarshal()` | yes | no |
| `json.Decoder` | no | yes |

When we encode data:

|  | `[]byte` | `io.Writer` |
|:--|:--|:--|
| `json.Marshal()` | yes | no |
| `json.Encoder` | no | yes |

That's a general rule. You may ask why? The answer to the question is developer experience. Let's consider an example where we have to read a body from a request. Let's use both `json.Unmarshal()` and `json.Decoder`. The `Request.Body` implements `io.Reader` interface so we can use this fact.

```go
req := CreateOrderRequest{}
if err := json.Decoder(r.Body).Decode(&req); err != nil {
  // handle the error
}

// the req is ready to use
```

We can write a similar program but use `json.Unmarshal()` to compare which code is more readable for us.

```go
req := CreateOrderRequest{}
body, err := io.ReadAll(r.Body)
if err != nil {
  // handle the error
}

if err = json.Unmarshal(body, &req); err != nil {
  // handle the error
}
```

There's more one difference that may tell us which one we should use. We can call `json.Decoder` and `json.Encoder` on a single `io.Reader` and `io.Writer` multiple times. It means that if the stream that we pass to the decoder contains multiple JSONs, we can create the decoder once but run `Decode()` multiple times.

```go
req := CreateOrderRequest{}
decoder := json.Decoder(r.Body)

for err := decoder.Decode(&req); err != nil {
	// handle single request
}
```

If you'd like to use `json.Decoder` but you have the `[]byte` you can wrap it with a buffer and use it instead.

```go
var body []byte
buf := bytes.NewBuffer(body)

decoder := json.Decoder(buf)
for err := decoder.Decode(&req); err != nil {
	// handle single request
}
```

### The performance?

I wrote simple benchmarks to compare both approaches.

```go
package jsons

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
)

var j = []byte(`{"user":"Johny Bravo","items":[{"id":"4983264583302173928","qty": 5}]}`)
var createRequest = CreateOrderRequest{
	User: "Johny Bravo",
	Items: []OrderItem{
		{ID: "4983264583302173928", Qty: 5},
	},
}
var err error
var body []byte

type OrderItem struct {
	ID  string `json:"id"`
	Qty int    `json:"qty"`
}

type CreateOrderRequest struct {
	User  string      `json:"user"`
	Items []OrderItem `json:"items"`
}

func BenchmarkJsonUnmarshal(b *testing.B) {
	b.ReportAllocs()
	req := CreateOrderRequest{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = json.Unmarshal(j, &req)
	}
}

func BenchmarkJsonDecoder(b *testing.B) {
	b.ReportAllocs()
	req := CreateOrderRequest{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		buff := bytes.NewBuffer(j)
		b.StartTimer()

		decoder := json.NewDecoder(buff)
		err = decoder.Decode(&req)
	}
}

func BenchmarkJsonMarshal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		body, err = json.Marshal(createRequest)
	}
}

func BenchmarkJsonEncoder(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		encoder := json.NewEncoder(io.Discard)
		err = encoder.Encode(createRequest)
	}
}
```

When we run them we'll see that `json.Unmarshal()` is about 3 times faster than `json.Decoder`. On the other hand, both `json.Marshal()` and `json.Encoder` have similar performance. At least with the input data I prepared.

```
BenchmarkJsonUnmarshal-10        1345796               894.4 ns/op           336 B/op          9 allocs/op
BenchmarkJsonDecoder-10           522276              2226 ns/op            1080 B/op         13 allocs/op
BenchmarkJsonMarshal-10          6257662               193.1 ns/op           128 B/op          2 allocs/op
BenchmarkJsonEncoder-10          6867033               174.9 ns/op            48 B/op          1 allocs/op
```

I encourage you to not take these or any other benchmarks as a go/no-go. You have to make similar tests in your application and then see if changing the function we use has any significant impact on the performance. Context is the king.

### Indenting
You probably saw that the JSON file produced by both `json.Marshal` or `json.Encoder` is compacted. Meaning, it has no extra white spaces that'd make it more human-readable. There's an alternative function called `json.MarshalIndent` that will help you format the output.

```go
	data := map[string]int{
		"a": 1,
		"b": 2,
	}

	b, err := json.MarshalIndent(data, "<prefix>", "<indent>")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
	
	// the output will be
	{
<prefix><indent>"a": 1,
<prefix><indent>"b": 2
<prefix>}
````

We can use the prefix to embed the new JSON into an already existing one and keep proper nesting.

## `MarshalJSON` and `UnmarshalJSON`
We can decide how a specific part of the JSON will be processed. We can achieve that we have to implement specific interfaces.

To be able to change the way our object is processed we have to implement one of those interfaces.

```go
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}
```

### `UnmarshalJSON` example
It will be easier to explain how it works with an example. Let's say we receive JSON with the PC's specification. The problem is that we receive the RAM size and storage size in bytes but we need it in a more human-readable format.

```json
{
	"cpu": "Intel Core i5",
	"operatingSystem": "Windows 11",
	"memory": 17179869184,
	"storage": 274877906944
}
```

Pretty unreadable, isn't it? Let's prepare our struct that will store these data.

```go
type PC struct {
	CPU string
	OperatingSystem string
	Memory string
	Storage string
}
```

To handle this case correctly we have to introduce a new type that will be an alias to a string. We'll implement the `UnmarhalJSON` method for it.

```go
type Memory string

func (m *Memory) UnmarshalJSON(b []byte) error {
	size, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	for i, d := range memorySizes {
		if size > d {
			*m = Memory(fmt.Sprintf("%d %s", size/d, sizeSuffixes[i]))
			return nil
		}
	}

	*m = Memory(fmt.Sprintf("%d b", size))
	return nil
}
```

We convert `[]byte` to an integer and then calculate the size in a human-readable format. The full source code is available here: https://goplay.tools/snippet/UfszC3iDvZW.

## FAQ

### What if I don't know the schema?
If you're not sure about the whole schema or part of it you have some options to handle it. One of the ways to go is using maps. Let's say we'll receive a JSON but we want to process it dynamically.

```go
req := map[string]interface{}{}
if err != json.Decoder(r.Body).Decode(&req); err != nil {
  // handle the err
}
```

We put the whole data into the map. Now, we can iterate over it and put our custom logic there. We'll need to use the `reflect` package to determine the type of value.

```go
	for k, v := range req {
		refVal := reflect.TypeOf(v)
		fmt.Printf("the key '%s' contains the value of type %s\n", k, refVal)
	}
	
	/* sample output:
	the key 'two' contains the value of type string
	the key 'three' contains the value of type float64
	the key 'one' contains the value of type int
	*/
```


### I cannot see my fields in JSON after marshaling
It can be caused by two things:
* the field isn't public (doesn't start with a capital letter)
* it's marked to be ignored using the struct tag: `json:"-"`

### Can I skip the error check in `Marshal()` method?
The general answer is `NO` but... I sometimes do it :)
If you can cover unsuccessful marshaling in your unit tests, I think it's OK to do it. Please just remember about adding a comment that it's a handler somewhere else.

On the other hand, is it worth making things a bit more complicated just to save one `if` statement? I'm not sure about it. It has chance to be an [unpopular opinion](https://changelog.com/gotime).

### If the std `json` package good enough?
I'd say 99% of the answer is `YES`. If you process huge JSON files or a lot of them and it's a significant part of the work, you may start seeking some alternatives. Otherwise, I think it won't disappoint you.

## Outside of the standard library
If you're looking for a faster alternative you can take a look at https://github.com/goccy/go-json. It's a drop-in replacement for the standard `encoding/json` package.

If the JSON is huge but you need only part of it, you can take a look at https://github.com/buger/jsonparser which allows you to just parse part of the whole file.

## Summary

I tried to cover everything that's needed to work with JSON in Go. If you have any other questions, feel free to use the comments section below. I'll be happy to answer any of them.

[^complex]: yes, Go has support for [complex numbers](https://go.dev/ref/spec#Complex_numbers). Only a few use it but I don't think it will be [removed from the language](https://github.com/golang/go/issues/19921).
