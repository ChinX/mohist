package nrouter

import (
	"bytes"
	"net/http"
	"net/url"
)

var delimiter = []byte("/")

const (
	levelOptimum patternType = iota
	levelVariable
	levelOptional
)

type (
	patternType int8
	Handle      func(rw http.ResponseWriter, req *http.Request, params *url.Values)
)

type Node struct {
	pattern  []byte
	level    patternType
	method   Handle
	children trunks
}

func NewNode() *Node {
	return &Node{children: newTrunks()}
}

func (n *Node) Add(path []byte, handler Handle) {
	path = bytes.Trim(path, string(delimiter))
	index := bytes.Index(path, delimiter)
	if index != -1 {
		n.pattern = path[:index]
		n.children.matchOrBuild(path[index+1:], handler)
	} else {
		if n.method != nil {
			panic("Ending point handler must be only")
		}
		n.method = handler
		n.pattern = path
	}
	n.level = ranking(path, index)
}

func (n *Node) checkEnding(path []byte, params *url.Values) bool {
	if n.level != levelOptimum {
		pattern := make([]byte,len(n.pattern)-int(n.level))
		copy(pattern,n.pattern)
		params.Add(string(pattern[int(n.level):]), string(path))
		return true
	}
	return bytes.Equal(n.pattern, path)
}

func (n *Node) matchEnding(path []byte, params *url.Values) Handle {
	if !n.checkEnding(path, params) {
		return nil
	}
	if n.method != nil {
		return n.method
	}
	if n.level != levelOptional && len(path) != 0 {
		return n.children.match(levelOptional, []byte(""), params)
	}
	return nil
}

func (n *Node) matchSelf(level patternType, path []byte, params *url.Values) Handle {
	index := bytes.Index(path, delimiter)
	if index == -1 {
		return n.matchEnding(path, params)
	}
	if n.checkEnding(path[:index], params) {
		return n.children.match(level, path[index+1:], params)
	}
	return nil
}

func (n *Node) Match(path []byte, params *url.Values) Handle {
	path = bytes.Trim(path, string(delimiter))
	return n.matchSelf(levelOptimum, path, params)
}

type trunks [][]*Node

func newTrunks() trunks {
	return make([][]*Node, 3)
}

func (t trunks) matchOrBuild(path []byte, handler Handle) {
	target := path
	index := bytes.Index(path, delimiter)
	if index != -1 {
		target = path[:index]
	}

	var aNode *Node
	for i := 0; i < len(t); i++ {
		children := t[i]
		for _, child := range children {
			if child != nil && bytes.Equal(child.pattern, target) {
				aNode = child
				break
			}
		}
	}

	if aNode == nil {
		aNode = NewNode()
		num := int(ranking(target, index))
		t[num] = append(t[num], aNode)
	}
	aNode.Add(path, handler)
}

func (t trunks) match(level patternType, path []byte, params *url.Values) Handle {
	for i := 0; i < len(t); i++ {
		children := t[i]
		for _, child := range children {
			if handler := child.matchSelf(level, path, params); handler != nil {
				return handler
			}
		}
	}
	return nil
}

func ranking(path []byte, index int) (level patternType) {
	switch path[0] {
	case ':':
		level = levelVariable
	case '?':
		if index == -1 {
			level = levelOptional
		} else {
			panic("Optional pattern must be ending point")
		}
	default:
		level = levelOptimum
	}
	return
}
