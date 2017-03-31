package test

import (
	"log"
	"strings"
	"testing"
	"unsafe"
)

func BenchmarkTrimString(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		addToTrimString(str)
	}
}

func BenchmarkTrimByte(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		addToTrimByte(str)
	}
}

func BenchmarkTrimBytes(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		addToTrimBytes(str)
	}
}

func TestTrim(t *testing.T) {
	//log.Println(addToTrimString("pattern"))
	//log.Println(addToTrimByte("pattern"))
	//log.Println(addToTrimString("/pattern"))
	//log.Println(addToTrimByte("/pattern"))
	//log.Println(addToTrimString("pattern/"))
	log.Println(addToTrimByte("pattern/"))
	//log.Println(addToTrimString("/pattern/"))
	//log.Println(addToTrimByte("/pattern/"))
	//log.Println(addToTrimString("///pattern///"))
	//log.Println(addToTrimByte("///pattern///"))
}

func addToTrimString(pattern string) string {
	return "/" + strings.Trim(pattern, "/")
}

func addToTrimByte(pattern string) string {
	return "/" + TrimByte(pattern, '/')
}

func addToTrimBytes(pattern string) string {
	return "/" + TrimBytes(pattern, '/')
}

func TrimRightByte(s string, b byte) string {
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

func TrimLeftByte(s string, b byte) string {
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

func TrimByte(s string, b byte) string {
	ns := TrimLeftByte(s, b)
	if len(ns) == len(s) {
		return s
	}
	return TrimRightByte(ns, b)
}

func TrimBytes(pattern string, b byte) string {
	e := len(pattern)
	if e == 0 {
		return pattern
	}
	s := -1
	for i := 0; i < e; i++ {
		if pattern[i] != b {
			s = i
			break
		}
	}
	if s == -1 || s == e-1 {
		return pattern
	}
	for j := e - 1; j > s; j-- {
		if pattern[j] != b {
			e = j + 1
			break
		}
	}
	return pattern[s:e]
}

func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func BenchmarkPartPath(b *testing.B) {
	p := "/abc/cde/fgh/ijk/"
	for i := 0; i < b.N; i++ {
		PartPath(p)
	}
}

func BenchmarkPartPathFunc(b *testing.B) {
	p := "/abc/cde/fgh/ijk/"
	for i := 0; i < b.N; i++ {
		partPathFunc(p, func(part string) {

		})
	}
}

func BenchmarkPartPathFunc2(b *testing.B) {
	p := "/abc/cde/fgh/ijk/"
	for i := 0; i < b.N; i++ {
		partPathFunc(p, logTTT)
	}
}

func logTTT(part string)  {

}

func PartPath(p string) {
	s, e, l := 0, 0, len(p)
	for l != e+1 {
		s, e = partIn(p, s)
		s = e
	}
}

func TestPartPath(t *testing.T) {
	p := "//abc//cde//fgh//ijk//"
	partPathFunc(p, func(part string) {
		log.Println(part)
	})
}

func partPathFunc(path string, fn func(string)) {
	s, e, l := 0, 0, len(path)
	for l != e+1 {
		s, e = partIn(path, s)
		if e-1 > s {
			fn(path[s:e])
		}
		s = e
	}
}

func partFunc(path string, fn func(string, bool)) {
	s, e, l := 0, 0, len(path)
	ending := false
	for !ending {
		s, e = partIn(path, s)
		ending = (l == e+1)
		if e-1 > s {
			fn(path[s:e], ending)
		}
		s = e
	}
}

func partIn(path string, start int) (s, e int) {
	first := false
	s, e = start, start
	for ; e < len(path); e++ {
		if (path[e] == '/') == first {
			if first {
				e = e
				return
			} else {
				s = e
				first = true
			}
		}
	}
	e = e - 1
	return
}
