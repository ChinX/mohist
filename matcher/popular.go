package matcher

import "regexp"

var (
	normal = regexp.MustCompile(`^\w{6}\w*$`)
	phone  = regexp.MustCompile(`^1(3[0-9]|4[57]|5[0-35-9]|7[0135678]|8[0-9])\d{8}$`)
	email  = regexp.MustCompile(`^\w+(\.\w+)*@\w+(\.\w+)+$`)
)

// MatchPhone ...
func MatchPhone(v string) bool {
	if v == "" {
		return false
	}
	return phone.MatchString(v)
}

// MatchEmail ...
func MatchEmail(v string) bool {
	if v == "" {
		return false
	}
	return email.MatchString(v)
}

// MatchPassWord ...
func MatchPassWord(v string) bool {
	if v == "" {
		return false
	}
	return normal.MatchString(v)
}
