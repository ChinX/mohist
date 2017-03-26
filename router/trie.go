package router

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
)

type node struct {
	pattern  string
	level    patternType
	method   Handle
	children trunks
}

func newNode() *node {
	return &node{children: newTrunks()}
}

func (n *node) add(path string, handler Handle) {
	path = strings.Trim(path, "/")
	index := strings.Index(path, "/")
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

func (n *node) checkEnding(path string, params *url.Values) bool {
	if n.level != levelOptimum {
		params.Add(n.pattern[int(n.level):], path)
		return true
	}
	return n.pattern == path
}

func (n *node) matchEnding(path string, params *url.Values) Handle {
	if !n.checkEnding(path, params) {
		return nil
	}
	if n.method != nil {
		return n.method
	}
	if n.level != levelOptional && path != "" {
		return n.children.match(levelOptional, "", params)
	}
	return nil
}

func (n *node) matchSelf(level patternType, path string, params *url.Values) Handle {
	index := strings.Index(path, "/")
	if index == -1 {
		return n.matchEnding(path, params)
	}
	if n.checkEnding(path[:index], params) {
		return n.children.match(level, path[index+1:], params)
	}
	return nil
}

func (n *node) match(path string, params *url.Values) Handle {
	path = strings.Trim(path, "/")
	return n.matchSelf(levelOptimum, path, params)
}

type trunks [][]*node

func newTrunks() trunks {
	return make([][]*node, 3)
}

func (t trunks) matchOrBuild(path string, handler Handle) {
	target := path
	index := strings.Index(path, "/")
	if index != -1 {
		target = path[:index]
	}

	var aNode *node
	for i := 0; i < len(t); i++ {
		children := t[i]
		for _, child := range children {
			if child != nil && child.pattern == target {
				aNode = child
				break
			}
		}
	}

	if aNode == nil {
		aNode = newNode()
		num := int(ranking(target, index))
		t[num] = append(t[num], aNode)
	}
	aNode.add(path, handler)
}

func (t trunks) match(level patternType, path string, params *url.Values) Handle {
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

func ranking(path string, index int) (level patternType) {
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
