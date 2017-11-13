package matcher

import (
	"testing"
	"log"
	"encoding/json"
)

func TestTraverse(t *testing.T) {
	url := "///abc/def///ghi/jkl///"
	l := len(url)
	for part, next := "", 0; next < l; {
		part, next = Traverse(url, next)
		log.Println(part, next, l)
	}

	url = "abc"
	l = len(url)
	for part, next := "", 0; next < l; {
		part, next = Traverse(url, next)
		log.Println(part, next, l)
	}
}

func TestAddNode(t *testing.T) {
	addUrl := []string{
		"/accounts/",
		"/accounts/:account",
		"/accounts/:account/projects",
		"/accounts/:account/projects/:project",
		"/accounts/:account/projects/:project/files/?file",
		"/ccounts/",
		"/ccounts/:account",
		"/ccounts/:account/projects",
		"/ccounts/:account/projects/:project",
		"/ccounts/:account/projects/:project/files/*file",
	}
	node := newNode("root")
	for i := range addUrl{
		node.AddNode(addUrl[i], struct {}{})
	}
	byt, _ := json.Marshal(node)
	log.Println(string(byt))
}