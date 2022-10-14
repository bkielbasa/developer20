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
