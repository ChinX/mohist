package matcher

import (
	"testing"
	"log"
)

func TestTraverse(t *testing.T) {
	url := "///abc/def///ghi/jkl///"
	l := len(url)
	for part, next := "", 0; next < l; {
		part, next = Traverse(url, next)
		log.Println(part, next, l)
	}

	url = "///"
	l = len(url)
	for part, next := "", 0; next < l; {
		part, next = Traverse(url, next)
		log.Println(part, next, l)
	}
}

func TestAddNode(t *testing.T) {
	url := "///abc/:def///ghi/?jkl///"
	AddNode(url)
	url = "///:abc/def///:ghi/*jkl///"
	AddNode(url)
}