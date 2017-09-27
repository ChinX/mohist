package mux

import (
	"fmt"
	"regexp"
)

const (
	necessary  int8 = iota
	varied
	negligible
)

var (
	partMath, _ = regexp.Compile(`^(\?*:+)*[a-zA-Z]([a-zA-Z0-9]|[a-zA-Z0-9])*$`)
	matchErr    = fmt.Sprintf("Path part must be matched %s!", partMath.String())
)



type leaf struct {
	kind        int8
	pattern     string
	match       func(path string, params Params) bool
	//addToLeaves func(leaves []*leaf) ([]*leaf, bool)
}

func newLeaf(path string, ending bool) (l *leaf) {
	if !partMath.MatchString(path) {
		panic(matchErr)
	}
	switch path[0] {
	case '?':
		if !ending {
			panic("Negligible pattern must be ending point")
		}
		if len(path) <= 2 || path[1] != ':' {
			panic(matchErr)
		}
		l = &leaf{kind: negligible, pattern: path[2:]}
		l.match = l.matchVaried
		//l.addToLeaves = l.negligibleToLeaves
	case ':':
		if len(path) == 1 {
			panic(matchErr)
		}
		l = &leaf{kind: varied, pattern: path[1:]}
		l.match = l.matchVaried
		//l.addToLeaves = l.variedToLeaves
	default:
		l = &leaf{kind: necessary, pattern: path}
		l.match = l.matchNecessary
		//l.addToLeaves = l.necessaryToLeaves
	}
	return
}

func (l *leaf) equal(lf *leaf) bool {
	return lf.pattern == l.pattern
}

func (l *leaf) matchVaried(path string, params Params) bool {
	params.Set(l.pattern, path)
	return true
}

func (l *leaf) matchNecessary(path string, params Params) bool {
	return path == l.pattern
}

//func (l *leaf) negligibleToLeaves(leaves []*leaf) ([]*leaf, bool) {
//	return leaves, false
//}
//
//func (l *leaf) variedToLeaves(leaves []*leaf) ([]*leaf, bool) {
//	return append(leaves, l), true
//}
//
//func (l *leaf) necessaryToLeaves(leaves []*leaf) ([]*leaf, bool) {
//	count := len(leaves)
//	ls := make([]*leaf, count+1, 2*count+10)
//	ls[0] = l
//	copy(ls[1:], leaves)
//	return ls, true
//}

func traverseLeaf(path string, mark byte, start int) (n int, part *leaf, ending bool) {
	l := len(path)
	switch l {
	case start:
		n, ending = start, true
	case start + 1:
		if path[start] == mark {
			n, ending = start, true
		} else {
			n = l
			part = newLeaf(path, ending)
		}
	default:
		begin := false
		n = start
		for ; n < len(path); n++ {
			if (path[n] == mark) == begin {
				if begin {
					break
				} else {
					start = n
					begin = true
				}
			}
		}
		ending = bool(l == n)
		if n-1 > start {
			part = newLeaf(path[start:n], l == n)
		}
	}
	return
}
