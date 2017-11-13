package matcher

import "log"

var errorPath = "path \"%s\" in pattern \"%s\" must matched regex \"^[a-zA-Z](([a-zA-Z0-9_.]*[a-zA-Z0-9]+))*$\""

type Node struct {
	staticLen int
	Pattern   string
	method    interface{} //http.HandlerFunc
	Children  []*Node
	Elastic   *Node
	Wide      *Node
}

func IsAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func IsDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func IsUnderling(ch byte) bool {
	return ch == '_'
}

func IsDot(ch byte) bool {
	return ch == '.'
}

func IsAlnum(ch byte) bool {
	return IsAlpha(ch) || IsDigit(ch)
}

func checkPart(part string) bool {
	l := len(part)
	switch l {
	case 0:
		return false
	case 1:
		return IsAlpha(part[0])
	default:
		l = l - 1
		if !IsAlpha(part[0]) || !IsAlnum(part[l]) {
			return false
		}
		for i := 1; i < l; i++ {
			if !IsAlnum(part[i]) && !IsDot(part[i]) && !IsUnderling(part[i]) {
				return false
			}
		}
	}
	return true
}

func newNode(pattern string) *Node {
	return &Node{
		Pattern:  pattern,
		Children: make([]*Node, 0, 10),
	}
}

func (n *Node) AddNode(pattern string, handler interface{}) {
	target := n
	length := len(pattern)
	for path, next := "", 0; next < length; {
		path, next = Traverse(pattern, next)
		if pl := len(path); pl == 0 || pl == 1 && !IsAlpha(path[0]) {
			log.Panicf("path \"%s\" in pattern \"%s\" must be an Alpha string", path, pattern)
		}

		var nn *Node
		switch path[0] {
		case '?':
			target.Elastic = endpointNode(target.Elastic, pattern, path[1:], handler, next != length)
		case '*':
			target.Wide = endpointNode(target.Wide, pattern, path[1:], handler, next != length)
		case ':':
			path = path[1:]
			if !checkPart(path) {
				log.Panicf(errorPath, path, pattern)
			}
			for l := target.staticLen; l > 0; l-- {
				if target.Children[l-1].Pattern == path {
					nn = target.Children[l-1]
					break
				}
			}
			newn := childNode(nn, path, handler, next == length)
			if nn == nil {
				target.Children = append(target.Children, newn)
				target.staticLen += 1
				nn = newn
			}
		default:
			if !checkPart(path) {
				log.Panicf(errorPath, path, pattern)
			}

			for i, l := target.staticLen, len(target.Children); i < l; i++ {
				if target.Children[i].Pattern == path {
					nn = target.Children[i]
					break
				}
			}

			newn := childNode(nn, path, handler, next == length)
			if nn == nil {
				target.Children = append([]*Node{newn}, target.Children...)
				nn = newn
			}

		}
		target = nn
	}
}

func childNode(n *Node, path string, handler interface{}, ended bool) *Node {
	if n == nil {
		n = newNode(path)
	}
	if ended {
		if n.method != nil {
			panic("method of endpoint path mush be only")
		}
		n.method = handler
	}
	return n
}

func endpointNode(n *Node, pattern, path string, handler interface{}, ended bool) *Node {
	if ended {
		log.Panicf("path \"%s\" in pattern \"%s\" must be endpoint", path, pattern)
	}

	if !checkPart(path) {
		log.Panicf(errorPath, path, pattern)
	}

	if n != nil {
		log.Panicf("path \"%s\" in pattern \"%s\" mush be only", path, pattern)
	}

	return &Node{
		Pattern: path,
		method:  handler,
	}
}

func Traverse(pattern string, start int) (part string, next int) {
	from := -1
	next = len(pattern)
	for i := start; i < next; i++ {
		if from == -1 {
			if pattern[i] != '/' {
				from = i
			}
			continue
		}
		if part == "" {
			if pattern[i] == '/' {
				part = pattern[from:i]
			}
		} else if pattern[i] != '/' {
			next = i
			break
		}
	}
	if part == "" && from != -1 {
		part = pattern[from:next]
	}
	return
}
