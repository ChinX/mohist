package new

import (
	"net/http"
	"net/url"
	"strings"
)

const (
	levelOptimum patternType = iota
	levelVariable
	levelOptional
)

type (
	patternType int8
	Handle      func(w http.ResponseWriter, req *http.Request)
)

type Node struct {
	Pattern  string      `json:"pattern"`
	Level    patternType `json:"level"`
	Methods  string      `json:"methods"` //map[string]string
	Children []*Node     `json:"children"`
}

func (n *Node) add(method, path string, handler string) {
	path = strings.Trim(path, "/")
	target := path
	index := strings.Index(path, "/")
	if index != -1 {
		target = path[:index]
		n.build(method, path[index+1:], handler)
	} else {
		n.Methods = handler
	}
	n.Pattern = target
	switch target[0] {
	case ':':
		n.Level = levelVariable
	case '?':
		if index == -1 {
			n.Level = levelOptional
		} else {
			panic("Optional pattern must be ending point")
		}
	default:
		n.Level = levelOptimum
	}
}

func (n *Node) build(method, path string, handler string) {
	target := path
	index := strings.Index(path, "/")
	if index != -1 {
		target = path[:index]
	}
	var aNode *Node
	for _, child := range n.Children {
		if child != nil && child.Pattern == target {
			aNode = child
			break
		}
	}
	if aNode == nil {
		aNode = &Node{}
		n.Children = append(n.Children, aNode)
	}
	aNode.add(method, path, handler)
}

func (n *Node) checkEnding(method, path string, params *url.Values) bool {
	if n.Level != levelOptimum {
		params.Set(n.Pattern[int(n.Level):], path)
		return true
	}
	return n.Pattern == path
}

func (n *Node) matchEnding(method, path string, params *url.Values) string {
	if !n.checkEnding(method, path, params) {
		return ""
	}
	if handler := n.Methods; handler != "" {
		return handler
	}
	if n.Level != levelOptional && path != "" {
		return n.matchChildren(levelOptional, method, "", params)
	}
	return ""
}

func (n *Node) matchChildren(level patternType, method, path string, params *url.Values) string {
	for _, child := range n.Children {
		if child != nil {
			handler := child.matchSelf(level, method, path, params)
			if handler != "" {
				return handler
			}
		}
	}
	return ""
}

func (n *Node) matchSelf(level patternType, method, path string, params *url.Values) string {
	index := strings.Index(path, "/")
	if index == -1 {
		return n.matchEnding(method, path, params)
	}
	if n.checkEnding(method, path[:index], params) {
		return n.matchChildren(level, method, path[index+1:], params)
	}
	return ""
}

func (n *Node) match(method, path string, params *url.Values) string {
	path = strings.Trim(path, "/")
	return n.matchSelf(levelOptimum, method, path, params)
}
