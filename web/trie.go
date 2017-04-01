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
	var nn *Node
	// todo: 增加层级，更好的命中handler
	traverseFunc(path, func(part string, ending bool) {
		switch part[0] {
		case ':':
			if len(part) == 1 || !checkPart(part[1:]) {
				panic("Wide pattern must be param")
			}
			//todo:层级仅能有一个param handler
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
		case '?':
			if ending && len(part) > 2 && part[1] == ':' &&
				checkPart(part[2:]) && target.Wide == nil {
				nn = newNode()
				nn.Pattern = part[2:]
				nn.Method = handler
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
		}
		if ending {
			if nn.Method != nil {
				//if nn.Method != "" {
				panic("Ending point method mush be only")
			}
			nn.Method = handler
		}
		target = nn
		nn = nil
	})
}

func (n *Node) match(path string) (handler Handle, values Params) {
	values = make(Params,0,128)
	target := n
	var nn *Node
	traverseFunc(path, func(part string, ending bool) {
		// match static handler
		for i := 0; i < len(target.Statics); i++ {
			if target.Statics[i].Pattern == part {
				nn = target.Statics[i]
				break
			}
		}

		if nn == nil {
			// match param handler
			for i := 0; i < len(target.Params); i++ {
				//todo:同级多个param，子路由无法被匹配出来
				if target.Params[i].Method != nil {
					//if target.Params[i].Method != "" {
					nn = target.Params[i]
					values = append(values, newParam(nn.Pattern, part))
					break
				}
			}
		}

		if nn == nil {
			// match wide handler
			if target.Wide != nil {
				values = append(values, newParam(target.Wide.Pattern, part))
				handler = target.Wide.Method
			}
		}

		if ending && nn != nil {
			if nn.Method != nil {
				//if nn.Method != "" {
				handler = nn.Method
			} else if nn.Wide != nil {
				handler = nn.Wide.Method
			}
		}
		target = nn
		nn = nil
	})
	return
}
