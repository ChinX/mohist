package web

import "net/http"

type Handle func(http.ResponseWriter, *http.Request, Params)

//type Handle string
type Params []*Param
type Param struct {
	Key   string
	Value string
}

func (ps Params) Get(key string) (string, bool) {
	for _, entry := range ps {
		if entry.Key == key {
			return entry.Value, true
		}
	}
	return "", false
}

func newParam(key, val string) *Param {
	return &Param{Key: key, Value: val}
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
	path = TrimByte(path, '/')
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
	return n.marchChildren(TrimByte(path, '/'), make(Params, 0, 128), 0)
}

func (n *node) marchChildren(path string, values Params, s int) (handler Handle, val Params) {
	val = values
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
				handler, val = nn.marchChildren(path, val, e)
			}
			return
		}
	}

	// match param handler
	for i := 0; i < len(n.params); i++ {
		nn := n.params[i]
		val = append(val, newParam(nn.pattern, part))
		if ending {
			if nn.method != nil {
				handler = nn.method
			} else if nn.wide != nil {
				handler = nn.wide.method
			}
		} else {
			handler, val = nn.marchChildren(path, val, e)
			if handler == nil {
				continue
			}
		}
		return
	}

	// match wide handler
	if ending && n.wide != nil {
		val = append(val, newParam(n.wide.pattern, part))
		handler = n.wide.method
	}
	return
}
