package web

import (
	"fmt"
	"net/url"
)

const (
	levelStatic = iota
	levelParam
	levelWide
)

type matcher struct {
	level   int
	pattern string
	handler Handle
	statics []*matcher
	params  []*matcher
	wide    *matcher
}

func newMatcher(pattern string, level int) *matcher {
	return &matcher{
		level:   level,
		pattern: pattern,
		statics: make([]*matcher, 0, 20),
		params:  make([]*matcher, 0, 20),
	}
}

func (m *matcher) match(pattern string, params *url.Values) Handle {
	if pattern[0] == '/' {
		pattern = pattern[1:]
	}
	l := len(pattern)
	if pattern[l-1] == '/' {
		pattern = pattern[:l-2]
	}
	l = len(m.pattern)
	switch m.level {
	case levelStatic:
		if l > len(pattern) {
			return nil
		}
		i := l - 1
		for ; i >= 0; i-- {
			if pattern[i] != pattern[i] {
				return nil
			}
		}
	case levelParam:
		var vars string
		for i := 0; i < l; i++ {
			if pattern[i] == '/' {
				vars = pattern[:i]
				pattern = pattern[i:]
			}
		}
		if len(pattern) == l {
			params.Add(m.pattern, pattern)
		} else {
			params.Add(m.pattern, vars)
		}
	default:
		for i := l - 1; i >= 0; i-- {
			if pattern[i] == '/' {
				return nil
			}
		}
		params.Add(m.pattern, pattern)
		return m.handler

	}
	if len(pattern) == l {
		if m.handler != nil {
			return m.handler
		}
		if m.wide != nil {
			return m.wide.handler
		}
		return nil
	}

	if pattern[l] == '/' {
		for _, s := range m.statics {
			if h := s.match(pattern[l+1:], params); h != nil {
				return h
			}
		}

		for _, p := range m.params {
			if h := p.match(pattern[l+1:], params); h != nil {
				return h
			}
		}

		if m.wide != nil {
			return m.wide.match(pattern[l+1:], params)
		}
	}
	return nil
}

func (m *matcher) add(pattern string, handler Handle) {
	if pattern[0] == '/' {
		pattern = pattern[1:]
	}
	l := len(pattern)
	if pattern[l-1] == '/' {
		pattern = pattern[:l-2]
	}
	vars := pattern
	target := m
	for {
		l = len(vars)
		for i := 0; i < l; i++ {
			if pattern[i] == '/' {
				vars = pattern[:i]
				pattern = pattern[i+1:]
				break
			}
		}

		if len(pattern) == l {
			pattern = ""
			l = 0
		}

		level := levelStatic
		var matches []*matcher
		switch vars[0] {
		case ':':
			vars = vars[1:]
			level = levelParam
			matches = target.params
		case '?':
			if vars[1] == ':' {
				if l != 0 || target.wide != nil {
					panic(fmt.Sprintf("Handler must be only in ending when wildcard is \"%s\"", vars))
				}
				target.wide = newMatcher(vars[2:], levelWide)
				return
			} else {
				panic(fmt.Sprintf("Invalid wildcard in string \"%s\"", vars))
			}
		default:
			matches = target.statics
		}

		var nm *matcher
		for _, p := range matches {
			if p.pattern == vars && p.level == level {
				nm = p
				break
			}
		}
		if nm == nil {
			nm = newMatcher(vars, level)
			target.params = append(target.params, nm)
		}
		if l == 0 {
			if nm.handler != nil {
				panic(fmt.Sprintf("Handler is exist when wildcard is \"%s\"", vars))
			}
			nm.handler = handler
			return
		}
		vars = pattern
		target = nm
	}
}
