package web

import "net/url"

const (
	levelStatic patternType = iota
	levelParam
	levelWide
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
	if path[0] == '/' {
		path = path[1:]
	}
	l := len(path)
	if path[l-1] == '/' {
		l -= 1
		path = path[:l]
	}
	i := 0
	for ; i < l; i++ {
		if path[i] == '/' {
			break
		}
	}
	if i != l {
		n.pattern = path[:i]
		n.children.matchOrBuild(path[i+1:], handler)
	} else {
		if n.method != nil {
			panic("Ending point handler must be only")
		}
		n.method = handler
		n.pattern = path
	}
	n.level = ranking(path, i == l)
}

func (n *node) checkEnding(path string, params *url.Values) bool {
	if n.level != levelStatic {
		params.Add(n.pattern[int(n.level):], path)
		return true
	}
	l := len(path)
	if len(n.pattern) == l {
		i := 0
		for ; i < l; i++ {
			if n.pattern[i] != path[i] {
				return false
			}
		}
	}
	return true
}

func (n *node) matchEnding(path string, params *url.Values) Handle {
	if !n.checkEnding(path, params) {
		return nil
	}
	if n.method != nil {
		return n.method
	}
	if n.level != levelWide && len(path) != 0 {
		return n.children.match(levelWide, "", params)
	}
	return nil
}

func (n *node) matchSelf(level patternType, path string, params *url.Values) Handle {
	l := len(path)
	i := 0
	for ; i < l; i++ {
		if path[i] == '/' {
			break
		}
	}
	if i == l {
		return n.matchEnding(path, params)
	}
	if n.checkEnding(path[:i], params) {
		return n.children.match(level, path[i+1:], params)
	}
	return nil
}

func (n *node) match(path string, params *url.Values) Handle {
	if path[0] == '/' {
		path = path[1:]
	}
	l := len(path) - 1
	if path[l] == '/' {
		path = path[:l]
	}
	return n.matchSelf(levelStatic, path, params)
}

type trunks [][]*node

func newTrunks() trunks {
	return make([][]*node, 3)
}

func (t trunks) matchOrBuild(path string, handler Handle) {
	leg := len(path)
	index := 0
	for ; index < leg; index++ {
		if path[index] == '/' {
			break
		}
	}
	var aNode *node
	for i := 0; i < len(t); i++ {
		children := t[i]
		for _, child := range children {
			if child != nil && len(child.pattern) == index {
				j := 0
				for ; j < index; j++ {
					if child.pattern[j] != path[j] {
						break
					}
				}
				if j == index {
					aNode = child
					break
				}
			}
		}
	}

	if aNode == nil {
		aNode = newNode()
		num := int(ranking(path[:index], index == leg))
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

func ranking(path string, ending bool) (level patternType) {
	switch path[0] {
	case ':':
		level = levelParam
	case '?':
		if ending {
			level = levelWide
		} else {
			panic("Optional pattern must be ending point")
		}
	default:
		level = levelStatic
	}
	return
}
