package internal

import (
	"log"
	"testing"
)

func TestUrlSplit(t *testing.T) {
	parse, err := UrlSplit("///accounts/account/users/user///")
	t.Log(parse, err)
}

func BenchmarkUrlSplit(t *testing.B) {
	for i := 0; i < t.N; i++ {
		UrlSplit("///accounts/account/users/user///")
	}
}

func TestTraverse(t *testing.T) {
	url := Trim("///accounts/account/users/user///", '/')
	part := ""
	var err error
	for start := 0; err == nil; {
		start, part, err = Traverse(url, start, '/')
		log.Println(start, part, err)
	}
}

func BenchmarkTraverse(t *testing.B) {
	url := Trim("///accounts/account/users/user///", '/')
	for i := 0; i < t.N; i++ {
		var err error
		for start := 0; err == nil; {
			start, _, err = Traverse(url, start, '/')
		}
	}
}
