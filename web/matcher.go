package web

import (
	"log"

	"github.com/chinx/mohist/validator"
)

const partRegex = "[a-z]([a-z0-9_]*[a-z0-9]+)$"

type Matcher struct {
	pattern string
	handler func(pattern string, part string) (key string, matched bool)
}

func (m *Matcher) Match(part string, params Params) bool {
	key, matched := m.handler(m.pattern, part)
	if !matched {
		return false
	}
	params.Set(key, part)
	return true
}

func Pattern(part string) (m *Matcher) {
	switch part[0] {
	case ':':
		if len(part) == 1 || checkPart(part[1:]) {
			log.Panicf("Param part must be \"^:%s\"", partRegex)
		}
		m = &Matcher{part[1:], paramMatch}
	case '?':
		if len(part) < 3 || part[1] != ':' || !checkPart(part[2:]) {
			log.Panicf("Wide part must be \"^?:%s\"", partRegex)
		}
		m = &Matcher{part[2:], paramMatch}
	default:
		if !checkPart(part) {
			log.Panicf("Static part must be \"^%s\"", partRegex)
		}
		m = &Matcher{part, staticMatch}
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
