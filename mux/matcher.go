package mux

import (
	"github.com/chinx/mohist/validator"
)

const (
	partRegex = `Pattern must be matched "^(\?*:){0,1}[a-z](_*[a-z0-9]+)*$"!`

	STATIC MatchKind = iota
	PARAM
	WIDE
)

type MatchKind int8

type Matcher struct {
	pattern string
	kind    MatchKind
	handler func(pattern string, part string) (key string, matched bool)
}

func (m *Matcher) Match(part string) (key string, matched bool) {
	return m.handler(m.pattern, part)
}

func (m *Matcher) Kind() MatchKind {
	return m.kind
}

func Pattern(part string) (m *Matcher) {
	switch part[0] {
	case ':':
		if len(part) == 1 || checkPart(part[1:]) {
			panic(partRegex)
		}
		m = &Matcher{part[1:], PARAM, paramMatch}
	case '?':
		if len(part) < 3 || part[1] != ':' || !checkPart(part[2:]) {
			panic(partRegex)
		}
		m = &Matcher{part[2:], WIDE, paramMatch}
	default:
		if !checkPart(part) {
			panic(partRegex)
		}
		m = &Matcher{part, STATIC, staticMatch}
	}
	return
}

func paramMatch(pattern string, part string) (key string, matched bool) {
	return pattern, part != ""
}

func staticMatch(pattern string, part string) (key string, matched bool) {
	return "", pattern == part
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
