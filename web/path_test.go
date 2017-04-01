package web

import (
	"regexp"
	"strings"
	"testing"
)

func BenchmarkTrimString(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		strings.Trim(str, "/")
	}
}

func BenchmarkTrimLeftByte(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		TrimLeftByte(str, '/')
	}
}

func BenchmarkTrimLeftString(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		strings.TrimLeft(str, "/")
	}
}

func BenchmarkTrimRightByte(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		TrimRightByte(str, '/')
	}
}

func BenchmarkTrimRightString(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		strings.TrimRight(str, "/")
	}
}

func BenchmarkTrimByte(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		TrimByte(str, '/')
	}
}

func BenchmarkConvertStrAndBytes(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		byt := []byte(str)
		str = string(byt)
	}
}

func BenchmarkConvertStrAndBytes2(b *testing.B) {
	str := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		byt := Str2bytes(str)
		str = Bytes2str(byt)
	}
}

func BenchmarkTravers(b *testing.B) {
	s := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		arr := strings.Split(s, "/")
		for j := 0; j < len(arr); j++ {
			if arr[j] != ""{}
		}

	}
}

func BenchmarkTraverseFunc(b *testing.B) {
	s := "////////////pattern////////////"
	for i := 0; i < b.N; i++ {
		traverseFunc(s, func(part string, ending bool) {
		})
	}

}

func BenchmarkCheckPart(b *testing.B) {
	s := "pattent"
	r := regexp.MustCompile(`^[A-Za-z_]*[A-Za-z0-9_.]*[A-Za-z0-9]$`)
	for i := 0; i < b.N; i++ {
		if r.MatchString(s) {
		}
	}
}

func BenchmarkCheckPart2(b *testing.B) {
	s := "pattent"
	for i := 0; i < b.N; i++ {
		if checkPart(s) {
		}
	}

}

func BenchmarkIsAlpha(b *testing.B) {
	s := "pattent"
	r := regexp.MustCompile(`^[A-Za-z_]+$`)
	for i := 0; i < b.N; i++ {
		if r.MatchString(s) {
		}
	}
}

func BenchmarkIsAlpha2(b *testing.B) {
	s := "pattent"
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(s); j++ {
			if IsAlpha(s[j]) {
			}
		}

	}
}

func BenchmarkIsDigit(b *testing.B) {
	s := "pattent"
	r := regexp.MustCompile(`^\d+$`)
	for i := 0; i < b.N; i++ {
		if r.MatchString(s) {
		}
	}
}

func BenchmarkIsDigit2(b *testing.B) {
	s := "pattent"
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(s); j++ {
			if IsDigit(s[j]) {
			}
		}

	}
}

func BenchmarkIsDot(b *testing.B) {
	s := "pattent"
	r := regexp.MustCompile(`^\.+$`)
	for i := 0; i < b.N; i++ {
		if r.MatchString(s) {
		}
	}
}

func BenchmarkIsDot2(b *testing.B) {
	s := "pattent"
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(s); j++ {
			if IsDot(s[j]) {
			}
		}

	}
}

func BenchmarkIsAlnum(b *testing.B) {
	s := "pattent"
	r := regexp.MustCompile(`^[A-Za-z0-9]+$`)
	for i := 0; i < b.N; i++ {
		if r.MatchString(s) {
		}
	}
}

func BenchmarkIsAlnum2(b *testing.B) {
	s := "pattent"
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(s); j++ {
			if IsAlnum(s[j]) {
			}
		}

	}
}
