package web

import (
	"net/http"

	"github.com/chinx/mohist/internal"
)

type Handle func(http.ResponseWriter, *http.Request, Params)

type node struct {
	pattern     string
	method      Handle
	paramHandle bool
	statics     []*node
	params      []*node
	wide        *node
}

func newNode() *node {
	return &node{
		statics: make([]*node, 0, 10),
		params:  make([]*node, 0, 10),
	}
}

func (n *node) addNode(path string, handler Handle) {
	target := n
	var nn *node
	// todo: 增加层级，更好的命中handler
	path = internal.Trim(path, '/')
	part, s, ending := "", 0, false
	for !ending {
		part, s, ending = traversePart(path, '/', s)
		switch part[0] {
		case ':':
			if len(part) == 1 || !checkPart(part[1:]) {
				panic("Wide pattern must be param")
			}
			if ending && target.paramHandle {
				panic("Ending point method mush be only")
			}
			for i := 0; i < len(target.params); i++ {
				if target.params[i].pattern == part[1:] {
					nn = target.params[i]
					break
				}
			}
			if nn == nil {
				nn = newNode()
				target.params = append(target.params, nn)
			}
			nn.pattern = part[1:]
			if ending {
				if nn.method != nil {
					panic("Ending point method mush be only")
				}
				nn.method = handler
				target.paramHandle = true
			}
		case '?':
			if !ending || len(part) < 3 || part[1] != ':' &&
				!checkPart(part[2:]) || target.wide != nil {
				panic("Wide pattern must be only in ending point")
			}
			nn = newNode()
			nn.pattern = part[2:]
			nn.method = handler
			target.wide = nn
		default:
			if !checkPart(part) {
				panic("Wide pattern must be static")
			}
			for i := 0; i < len(target.statics); i++ {
				if target.statics[i].pattern == part {
					nn = target.statics[i]
					break
				}
			}
			if nn == nil {
				nn = newNode()
				target.statics = append(target.statics, nn)
			}
			nn.pattern = part
			if ending {
				if nn.method != nil {
					panic("Ending point method mush be only")
				}
				nn.method = handler
			}
		}
		target = nn
		nn = nil
	}
}

func (n *node) match(path string) (Handle, Params) {
	values := NewParams()
	return n.marchChildren(internal.Trim(path, '/'), values, 0), values
}

func (n *node) marchChildren(path string, values Params, s int) (handler Handle) {
	part, e, ending := traversePart(path, '/', s)

	// match static handler
	for i := 0; i < len(n.statics); i++ {
		if n.statics[i].pattern == part {
			nn := n.statics[i]
			if ending {
				if nn.method != nil {
					handler = nn.method
				} else if nn.wide != nil {
					handler = nn.wide.method
				}
			} else {
				handler = nn.marchChildren(path, values, e)
			}
			return
		}
	}

	// match param handler
	for i := 0; i < len(n.params); i++ {
		nn := n.params[i]
		values.Set(nn.pattern, part)
		if ending {
			if nn.method != nil {
				handler = nn.method
			} else if nn.wide != nil {
				handler = nn.wide.method
			}
		} else {
			handler = nn.marchChildren(path, values, e)
			if handler == nil {
				continue
			}
		}
		return
	}

	// match wide handler
	if ending && n.wide != nil {
		values.Set(n.wide.pattern, part)
		handler = n.wide.method
	}
	return
}

func traversePart(path string, b byte, s int) (part string, n int, ending bool) {
	l := len(path)
	switch l {
	case s:
		n, ending = s, true
	case s + 1:
		if path[s] == b {
			n, ending = s, true
		} else {
			part, n, ending = path, l, true
		}
	default:
		begin := false
		n = s
		for ; n < len(path); n++ {
			if (path[n] == b) == begin {
				if begin {
					break
				} else {
					s = n
					begin = true
				}
			}
		}
		ending = (l == n)
		if n-1 > s {
			part = path[s:n]
		}
	}
	return
}
