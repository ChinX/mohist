package web

import (
	"net/url"
	"strings"
)

type route struct {
	handler Handle
}

type node struct {
	children     []*node
	component    string
	isNamedParam bool
	methods      map[string]*route
}

func (n *node) addNode(method, path string, handler Handle) {
	components := strings.Split(strings.ToLower(path), "/")[1:]
	count := len(components)
	for count > 0 {
		aNode, component := n.traverse(components, nil)
		if aNode.component == component && count == 1 {
			r := &route{handler}
			aNode.methods[method] = r
			return
		}
		newNode := node{component: component, isNamedParam: false, methods: make(map[string]*route)}
		if len(component) > 0 && component[0] == ':' {
			newNode.isNamedParam = true
		}
		if count == 1 {
			r := &route{handler}
			aNode.methods[method] = r
		}
		aNode.children = append(aNode.children, &newNode)
		count--
	}
}

func (n *node) traverse(components []string, params *url.Values) (*node, string) {
	component := components[0]
	if len(n.children) > 0 {
		for _, child := range n.children {
			if component == child.component || child.isNamedParam {
				if child.isNamedParam && params != nil {
					params.Set(child.component[1:], component)
				}
				next := components[1:]
				if len(next) > 0 {
					return child.traverse(next, params)
				}
				return child, component
			}
		}
	}
	return n, component
}
