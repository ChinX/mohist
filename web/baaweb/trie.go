package baaweb

import (
	"net/http"
	"net/url"
)

type Handle func(http.ResponseWriter, *http.Request, *url.Values)

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
	traverseFunc(path, func(part string, ending bool) {
		switch part[0] {
		case ':':
			if len(part) == 1 || !checkPart(part[1:]) {
				panic("Wide pattern must be param")
			}
			//todo:增加同层param的判断
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
				panic("Ending point method mush be only")
			}
			nn.Method = handler
		}
		target = nn
		nn = nil
	})
}

func (n *Node) match(path string, values *url.Values) (handler Handle) {
	target := n
	var nn *Node
	traverseFunc(path, func(part string, ending bool) {
		for i := 0; i < len(target.Statics); i++ {
			if target.Statics[i].Pattern == part {
				nn = target.Statics[i]
				break
			}
		}
		if nn == nil {
			for i := 0; i < len(target.Params); i++ {
				if target.Params[i].Pattern == part {
					nn = target.Params[i]
					values.Add(nn.Pattern, part)
					break
				}
			}
			if nn == nil {
				nn = target.Wide
			}
			if nn != nil {

			}
		}
		if nn != nil {
			if handler = nn.Method; handler == nil {
				if nn.Wide != nil {
					handler = nn.Wide.Method
				}
			}
		}
		target = nn
		nn = nil
	})
	return
}
