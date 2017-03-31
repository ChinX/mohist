package web

import "net/url"

type Handle string

type Node struct {
	Pattern string
	Method  Handle
	Statics []*Node
	Params  []*Node
	Wide    *Node
}

func newNode() *Node {
	return &Node{
		Statics: make([]*Node, 0, 10),
		Params:  make([]*Node, 0, 10),
	}
}

func (n *Node) addNode(path string, handler Handle) {
	target := n
	partFunc(path, func(part string, ending bool) {
		target = target.addChild(part, ending, handler)
	})
}

func (n *Node) addChild(part string, ending bool, handler Handle) (nn *Node) {
	switch part[0] {
	case ':':
		if len(part) == 1 || !checkPart(part[1:]) {
			panic("Wide pattern must be param")
		}
		//todo:增加同层param的判断
		for i := 0; i < len(n.Params); i++ {
			if n.Params[i].Pattern == part[1:] {
				nn = n.Params[i]
				break
			}
		}
		if nn == nil {
			nn = newNode()
			n.Params = append(n.Params, nn)
		}
		nn.Pattern = part[1:]
	case '?':
		if ending && len(part) > 2 && part[1] == ':' &&
			checkPart(part[2:]) && n.Wide == nil {
			nn = newNode()
			nn.Pattern = part[2:]
			nn.Method = handler
			n.Wide = nn
			return
		}
		panic("Wide pattern must be only in ending point")
	default:
		if !checkPart(part) {
			panic("Wide pattern must be static")
		}
		for i := 0; i < len(n.Statics); i++ {
			if n.Statics[i].Pattern == part {
				nn = n.Statics[i]
				break
			}
		}
		if nn == nil {
			nn = newNode()
			n.Statics = append(n.Statics, nn)
		}
		nn.Pattern = part
	}
	if ending {
		if nn.Method != "" {
			panic("Ending point method mush be only")
		}
		nn.Method = handler
	}
	return
}

func (n *Node)match(path string, values *url.Values) (handler Handle) {
	target := n
	partFunc(path, func(part string, ending bool) {
		target, handler = target.matchChild(part, ending, values)
	})
	return
}

func (n *Node)matchChild(part string, ending bool, values*url.Values) (nn *Node, handler Handle) {
	for i := 0; i < len(n.Statics); i++ {
		if n.Statics[i].Pattern == part{
			nn = n.Statics[i]
			break
		}
	}
	if nn == nil{
		for i := 0; i < len(n.Params); i++ {
			if n.Params[i].Pattern == part{
				nn = n.Params[i]
				break
			}
		}
		if nn == nil{
			nn = n.Wide
		}
	}
	if nn == nil{
		return
	}
	if handler = nn.Method; handler == ""{
		if nn.Wide != nil{
			handler = nn.Wide.Method
		}
	}
	return
}

func partFunc(path string, fn func(string, bool)) {
	s, e, l := 0, 0, len(path)
	ending := false
	for !ending {
		s, e = partIn(path, s)
		ending = (l == e)
		if e-1 > s {
			fn(path[s:e], ending)
		}
		s = e
	}
}

func partIn(path string, start int) (s, e int) {
	first := false
	s, e = start, start
	for ; e < len(path); e++ {
		if (path[e] == '/') == first {
			if first {
				return
			} else {
				s = e
				first = true
			}
		}
	}
	return
}

func checkPart(part string) bool {
	for i := 0; i < len(part); i++ {
		if !isAlpha(part[i]) {
			return false
		}
	}
	return true
}

func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isAlnum(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}

func trimPath(p string) string {
	e := len(p)
	if e == 0 {
		return p
	}
	s := -1
	for i := 0; i < e; i++ {
		if p[i] != '/' {
			s = i
			break
		}
	}
	if s == -1 || s == e-1 {
		return p
	}
	for j := e - 1; j > s; j-- {
		if p[j] != '/' {
			e = j + 1
			break
		}
	}
	return p[s:e]
}
