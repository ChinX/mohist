package binary

import (
	"strings"
	"testing"
	"reflect"
	"unsafe"
)

func BenchmarkConvertStrAndBytes(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		byt := StrToBytes(str)
		str = BytesToStr(byt)
	}
}

func BenchmarkConvertStrAndBytes2(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		byt := []byte(str)
		str = string(byt)
	}
}

func BenchmarkTrim(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		Trim(str, '/')
	}
}

func BenchmarkTrimStr(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		strings.Trim(str, "/")
	}
}

func BenchmarkTrimLeft(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		TrimLeft(str, '/')
	}
}

func BenchmarkTrimLeftStr(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		strings.TrimLeft(str, "/")
	}
}

func BenchmarkTrimRight(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		TrimRight(str, '/')
	}
}

func BenchmarkTrimRightStr(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		strings.TrimRight(str, "/")
	}
}

type MyStruct struct {
	A string `json:"a"`
	B string `json:"b"`
}

func TestStructToBytes(t *testing.T) {
	s := MyStruct{A:"abc",B:"cde"}

	l := int(unsafe.Sizeof(MyStruct{}))
	t.Log(l)
	x := reflect.SliceHeader{
		Len: l,
		Cap: l,
		Data: uintptr(unsafe.Pointer(&s)),
	}
	t.Log(*(*[]byte)(unsafe.Pointer(&x)))

	t.Log(StructToBytes(s))
	t.Log(StructToBytes(&s))
}
