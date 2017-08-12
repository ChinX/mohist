package internal

func TrimRight(s string, b byte) string {
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

func TrimLeft(s string, b byte) string {
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

func Trim(s string, b byte) string {
	ns := TrimLeft(s, b)
	if len(ns) == len(s) {
		return s
	}
	return TrimRight(ns, b)
}