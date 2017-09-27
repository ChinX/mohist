package web

import (
	"io"
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
	var err error
	for next, part := 0, ""; err == nil; {
		next, part, err = internal.Traverse(path, next, '/')
		if err != nil && err != io.EOF {
			panic(err)
		}
		switch part[0] {
		case ':':
			if len(part) == 1 || !checkPart(part[1:]) {
				panic("Wide pattern must be param")
			}
			if err == io.EOF && target.paramHandle {
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
			if err == io.EOF {
				if nn.method != nil {
					panic("Ending point method mush be only")
				}
				nn.method = handler
				target.paramHandle = true
			}
		case '?':
			if err != io.EOF || len(part) < 3 || part[1] != ':' &&
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
			if err == io.EOF {
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

func (n *node) match(path string) (Handle, Params, error) {
	values := NewParams()
	handler, err := n.marchChildren(internal.Trim(path, '/'), values, 0)
	return handler, values, err
}

func (n *node) marchChildren(path string, values Params, start int) (handler Handle, err error) {
	next, part := start, ""
	next, part, err = internal.Traverse(path, next, '/')
	if err != nil && err != io.EOF {
		return
	}

	// match static handler
	for i := 0; i < len(n.statics); i++ {
		if n.statics[i].pattern == part {
			nn := n.statics[i]
			if err == io.EOF {
				if nn.method != nil {
					handler = nn.method
				} else if nn.wide != nil {
					handler = nn.wide.method
				}
			} else {
				handler, err = nn.marchChildren(path, values, next)
			}
			return
		}
	}

	// match param handler
	for i := 0; i < len(n.params); i++ {
		nn := n.params[i]
		values.Set(nn.pattern, part)
		if err == io.EOF {
			if nn.method != nil {
				handler = nn.method
			} else if nn.wide != nil {
				handler = nn.wide.method
			}
		} else {
			handler, err = nn.marchChildren(path, values, next)
			if err != nil && err != io.EOF {
				return
			}
			if handler == nil {
				continue
			}
		}
		return
	}

	// match wide handler
	if err == io.EOF && n.wide != nil {
		values.Set(n.wide.pattern, part)
		handler = n.wide.method
	}
	return
}
