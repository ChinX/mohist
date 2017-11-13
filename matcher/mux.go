package matcher

import (
	"log"
)

const (
	static = iota
	dynamic
	elastic
	wide
)

func IsAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func IsDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func IsUnderling(ch byte) bool {
	return ch == '_'
}

func IsDot(ch byte) bool {
	return ch == '.'
}

func IsAlnum(ch byte) bool {
	return IsAlpha(ch) || IsDigit(ch)
}

func AddNode(pattern string) {
	l := len(pattern)
	for path, next := "", 0; next < l; {
		path, next = Traverse(pattern, next)
		pl := len(path)
		kind := static
		switch pl {
		case 0:
			log.Panicf("path in pattern \"%s\" be empty", pattern)
		case 1:
			if !IsAlpha(path[0]) {
				log.Panicf("path \"%s\" in pattern \"%s\" not Alpha", path, pattern)
			}
		default:
			switch path[0] {
			case ':':
				kind = dynamic
				path = path[1:]
			case '?':
				kind = elastic
				path = path[1:]
			case '*':
				kind = wide
				path = path[1:]
			}
		}
		if next != l && kind > dynamic {
			log.Panicf("elastic or wide path \"%s\" in pattern \"%s\" must be endpoint", path, pattern)
		}
		arr := []string{"static", "dynamic", "elastic", "wide"}
		log.Println(arr[kind], path)
	}
}

func Traverse(pattern string, start int) (part string, next int) {
	from := -1
	next = len(pattern)
	for i := start; i < next; i++ {
		if from == -1 {
			if pattern[i] != '/' {
				from = i
			}
			continue
		}
		if part == "" {
			if pattern[i] == '/' {
				part = pattern[from:i]
			}
		} else if pattern[i] != '/' {
			next = i
			break
		}
	}
	return
}
