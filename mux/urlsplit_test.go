package mux

import (
	"log"
	"testing"

	"github.com/chinx/mohist/internal"
)

func TestUrlSplit(t *testing.T) {
	parse, err := UrlSplit("/accounts//account///users////user", false)
	t.Log(parse, err)
}

func BenchmarkUrlSplit(t *testing.B) {
	for i := 0; i < t.N; i++ {
		UrlSplit("/accounts//account///users////user", false)
	}
}

func TestTraverse(t *testing.T) {
	parse := NewParse(internal.Trim("accounts//account///users////user///", '/'))
	part, ending := "", false
	for !ending {
		part, ending = parse.Traverse()
		log.Println(part)
	}
}

func BenchmarkTraverse(t *testing.B) {
	parse := NewParse(internal.Trim("/accounts//account///users////user", '/'))
	for i := 0; i < t.N; i++ {
		ending := false
		for !ending {
			_, ending = parse.Traverse()
		}
	}
}
