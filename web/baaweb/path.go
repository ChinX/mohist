package baaweb

import "unsafe"

func traverseFunc(path string, fn func(string, bool)) {
	if len(path) == 0 {
		fn(path, true)
	}
	path = TrimRightByte(path, '/')
	s, e, l := 0, 0, len(path)
	start, ending := false, (l == e)
	for !ending {
		for ; e < len(path); e++ {
			if (path[e] == '/') == start {
				if start {
					break
				} else {
					s = e
					start = true
				}
			}
		}
		ending = (l == e)
		if e-1 > s {
			fn(path[s:e], ending)
		}
		s = e
		start = false
	}
}

func checkPart(part string) bool {
	l := len(part)
	switch l {
	case 0:
		return false
	case 1:
		return IsAlpha(part[0])
	default:
		l = l - 1
		if !IsAlnum(part[l]) {
			return false
		}
		for i := 1; i < l; i++ {
			if !IsAlnum(part[i]) && !IsDot(part[i]) {
				return false
			}
		}
	}
	return true
}

func IsAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func IsDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func IsDot(ch byte) bool {
	return ch == '.'
}

func IsAlnum(ch byte) bool {
	return IsAlpha(ch) || IsDigit(ch)
}

func TrimRightByte(s string, b byte) string {
	if s == "" {
		return s
	}
	i := len(s)
	for ; i > 0; i-- {
		if s[i-1] != b {
			break
		}
	}
	return s[:i]
}

func TrimLeftByte(s string, b byte) string {
	if s == "" {
		return s
	}
	i := 0
	for ; i < len(s); i++ {
		if s[i] != b {
			break
		}
	}
	return s[i:]
}

func TrimByte(s string, b byte) string {
	ns := TrimLeftByte(s, b)
	if len(ns) == len(s) {
		return s
	}
	return TrimRightByte(ns, b)
}

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
