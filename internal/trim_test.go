package internal

import (
	"testing"
	"strings"
)

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
