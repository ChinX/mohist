package new

import (
	"net/url"
	"testing"
	"encoding/json"
)

func TestAdd(t *testing.T) {
	n := &Node{}
	n.add("", "/accounts/", "one")
	n.add("", "accounts/:account", "two")
	n.add("", "accounts/:account/projects", "three")
	n.add("", "accounts/:account/projects/:project", "four")
	n.add("", "accounts/:account/projects/:project/files/?:file", "six")

	byteArr, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(byteArr))


	param := &url.Values{}
	t.Log(n.match("", "/accounts/", param), param)
	param = &url.Values{}
	t.Log(n.match("", "accounts/account", param), param)
	param = &url.Values{}
	t.Log(n.match("", "accounts/abc/projects", param), param)
	param = &url.Values{}
	t.Log(n.match("", "acco/account/projects/project", param), param)
	param = &url.Values{}
	t.Log(n.match("", "accounts/account/projects/project/files", param), param)
	param = &url.Values{}
	t.Log(n.match("", "accounts/account/projects/project/files/file", param), param)
	param = &url.Values{}
	t.Log(n.match("", "accounts/account/projects/project/files/file/1", param), param)
}

func TestMatch(t *testing.T) {
	n := newNode(levelOptimum, "accounts", "one",
		newNode(levelVariable, ":account", "two",
			newNode(levelOptimum, "projects", "three",
				newNode(levelVariable, ":project", "four",
					newNode(levelOptimum, "files", "",
						newNode(levelOptional, "?:file", "six", nil))))))
	param := &url.Values{}
	t.Log(n.match("", "/accounts/", param), param)
	param = &url.Values{}
	t.Log(n.match("", "accounts/account", param), param)
	param = &url.Values{}
	t.Log(n.match("", "accounts/abc/projects", param), param)
	param = &url.Values{}
	t.Log(n.match("", "acco/account/projects/project", param), param)
	param = &url.Values{}
	t.Log(n.match("", "accounts/account/projects/project/files", param), param)
	param = &url.Values{}
	t.Log(n.match("", "accounts/account/projects/project/files/file", param), param)
	param = &url.Values{}
	t.Log(n.match("", "accounts/account/projects/project/files/file/1", param), param)
}

func newNode(level patternType, pattern, methods string, node *Node) *Node {
	nNode := &Node{
		Pattern:  pattern,
		Level:    level,
		Methods:  methods,
		Children: make([]*Node, 0, 10),
	}
	if node != nil{
		nNode.Children = append(nNode.Children, node)
	}
	return nNode
}
