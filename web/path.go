package web

import "unsafe"

func traversePart(path string, b byte, s int) (part string, n int, ending bool) {
	l := len(path)
	switch l {
	case s:
		n, ending = s, true
	case s + 1:
		if path[s] == b {
			n, ending = s, true
		} else {
			part, n, ending = path, l, true
		}
	default:
		begin := false
		n = s
		for ; n < len(path); n++ {
			if (path[n] == b) == begin {
				if begin {
					break
				} else {
					s = n
					begin = true
				}
			}
		}
		ending = (l == n)
		if n-1 > s {
			part = path[s:n]
		}
	}
	return
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
