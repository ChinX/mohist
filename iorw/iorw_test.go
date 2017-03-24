package iorw

import (
	"bytes"
	"testing"
)

type TestBody struct {
	A string
	B string
}

func TestJSONBuffer(t *testing.T) {
	a, err := JSONBuffer(&TestBody{A: "JSON", B: "Buffer"})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(a.Bytes()))

	b := &TestBody{}
	if err := ReadJSON(a, b); err != nil {
		t.Error(err)
		return
	}
	t.Log(b)
}

func TestWriteJSON(t *testing.T) {
	a := bytes.NewBuffer([]byte{})
	if err := WriteJSON(a, &TestBody{A: "Write", B: "JSON"}); err != nil {
		t.Error(err)
		return
	}
	t.Log(string(a.Bytes()))

	b := &TestBody{}
	if err := ReadJSON(a, b); err != nil {
		t.Error(err)
		return
	}
	t.Log(b)

}
