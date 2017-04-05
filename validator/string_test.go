package validator

import (
	"regexp"
	"testing"
)

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
