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

type Node struct {
	Pattern     string
	method      Handle
	ParamHandle bool
	Statics     []*Node
	Params      []*Node
	Wide        *Node
}

func newNode() *Node {
	return &Node{
		Statics: make([]*Node, 0, 10),
		Params:  make([]*Node, 0, 10),
	}
}

func (n *Node) addNode(path string, handler Handle) {
	target := n
	var nn *Node
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
			if ending && target.ParamHandle{
				panic("Ending point method mush be only")
			}
			for i := 0; i < len(target.Params); i++ {
				if target.Params[i].Pattern == part[1:] {
					nn = target.Params[i]
					break
				}
			}
			if nn == nil {
				nn = newNode()
				target.Params = append(target.Params, nn)
			}
			nn.Pattern = part[1:]
			if ending {
				if nn.method != nil {
					panic("Ending point method mush be only")
				}
				nn.method = handler
				target.ParamHandle = true
			}
		case '?':
			if ending && len(part) > 2 && part[1] == ':' &&
				checkPart(part[2:]) && target.Wide == nil {
				nn = newNode()
				nn.Pattern = part[2:]
				nn.method = handler
				target.Wide = nn
				return
			}
			panic("Wide pattern must be only in ending point")
		default:
			if !checkPart(part) {
				panic("Wide pattern must be static")
			}
			for i := 0; i < len(target.Statics); i++ {
				if target.Statics[i].Pattern == part {
					nn = target.Statics[i]
					break
				}
			}
			if nn == nil {
				nn = newNode()
				target.Statics = append(target.Statics, nn)
			}
			nn.Pattern = part
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

func (n *Node) match(path string) (Handle, Params) {
	return n.marchChildren(TrimByte(path, '/'), make(Params, 0, 128), 0)
}

func (n *Node) marchChildren(path string, values Params, s int) (handler Handle, val Params) {
	val = values
	part, e, ending := traversePart(path, '/', s)

	// match static handler
	for i := 0; i < len(n.Statics); i++ {
		if n.Statics[i].Pattern == part {
			nn := n.Statics[i]
			if ending {
				if nn.method != nil {
					handler = nn.method
				} else if nn.Wide != nil {
					handler = nn.Wide.method
				}
			} else {
				handler, val = nn.marchChildren(path, val, e)
			}
			return
		}
	}

	// match param handler
	for i := 0; i < len(n.Params); i++ {
		nn := n.Params[i]
		val = append(val, newParam(nn.Pattern, part))
		if ending {
			if nn.method != nil {
				handler = nn.method
			} else if nn.Wide != nil {
				handler = nn.Wide.method
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
	if ending && n.Wide != nil {
		val = append(val, newParam(n.Wide.Pattern, part))
		handler = n.Wide.method
	}
	return
}
