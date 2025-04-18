package rest

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBindJSON(t *testing.T) {
	data := []byte(`{"name":"John Doe","age":30}`)
	dataList := []byte(`[{"name":"John Doe","age":30}]`)
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	v := []Person{}
	want := []Person{
		{
			Name: "John Doe",
			Age:  30,
		},
	}

	if err := BindJSONList(bytes.NewReader(data), &v); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(v, want) {
		t.Fatalf("expected %v, got %v", want, v)
	}

	if err := BindJSONList(bytes.NewReader(dataList), &v); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(v, want) {
		t.Fatalf("expected %v, got %v", want, v)
	}
}
