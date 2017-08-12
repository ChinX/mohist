package web

import (
	"net/http"

	"fmt"

	"github.com/chinx/mohist/internal"
	"github.com/chinx/mohist/validator"
)

type Handle func(http.ResponseWriter, *http.Request, Params)

//type Handle string
type Params []*param

type param struct {
	Key   string
	Value string
}

const maxParam = 128

func (ps Params) Get(key string) (string, bool) {
	for _, entry := range ps {
		if entry.Key == key {
			return entry.Value, true
		}
	}
	return "", false
}

func (ps Params) Set(key, val string) error {
	if len(ps) >= maxParam {
		return fmt.Errorf("the params length must less than %d", maxParam)
	}
	for _, entry := range ps {
		if entry.Key == key {
			return fmt.Errorf("the key \"%s\" is exist in params", entry.Key)
		}
	}
	ps = append(ps, &param{Key: key, Value: val})
	return nil
}

func NewParams() Params {
	return Params(make([]*param, 0, maxParam))
}

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
			if ending && len(part) > 2 && part[1] == ':' &&
				checkPart(part[2:]) && target.wide == nil {
				nn = newNode()
				nn.pattern = part[2:]
				nn.method = handler
				target.wide = nn
				return
			}
			panic("Wide pattern must be only in ending point")
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

func checkPart(part string) bool {
	l := len(part)
	switch l {
	case 0:
		return false
	case 1:
		return validator.IsAlpha(part[0])
	default:
		l = l - 1
		if !validator.IsAlnum(part[l]) {
			return false
		}
		for i := 1; i < l; i++ {
			if !validator.IsAlnum(part[i]) && !validator.IsDot(part[i]) {
				return false
			}
		}
	}
	return true
}
