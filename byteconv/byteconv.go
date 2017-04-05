package byteconv

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

func StrToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

func TrimRight(s string, b byte) string {
	if s == "" {
		return s
	}
	i := len(s)
	for ; i > 0; i-- {
		if s[i-1] != b {
			break
		}
	}
	return s[:i]
}

func TrimLeft(s string, b byte) string {
	if s == "" {
		return s
	}
	i := 0
	for ; i < len(s); i++ {
		if s[i] != b {
			break
		}
	}
	return s[i:]
}

func Trim(s string, b byte) string {
	ns := TrimLeft(s, b)
	if len(ns) == len(s) {
		return s
	}
	return TrimRight(ns, b)
}
