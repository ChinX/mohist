package mux

import (
	"net/http"

	"io"

	"github.com/chinx/mohist/internal"
)

type Handle func(http.ResponseWriter, *http.Request, Params)

type node struct {
	*leaf
	method      Handle
	endHandle   bool
	children    []*node
	paramHandle bool
	statics     []*node
	params      []*node
	wide        *node
}

func newNode() *node {
	return &node{
		children: make([]*node, 0, 10),
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
		matcher := Pattern(part)
		switch matcher.Kind() {
		case STATIC:
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
		case PARAM:
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
		case WIDE:
			if err != io.EOF || target.wide != nil {
				panic("Wide pattern must be only in ending point")
			}
			nn = newNode()
			nn.pattern = part[2:]
			nn.method = handler
			target.wide = nn
		}
	}
	target = nn
	nn = nil
}
