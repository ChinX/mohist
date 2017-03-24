package random

import "crypto/rand"

var pool = []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// NewStr create a random string
func NewStr(n int) (string, error) {
	return randomStr(pool, n)
}

// UpperStr create a random upper string
func UpperStr(n int) (string, error) {
	return randomStr(pool[36:], n)
}

// NoUpperStr create a random string has not upper
func NoUpperStr(n int) (string, error) {
	return randomStr(pool[:36], n)
}

// LowerStr create a random lower string
func LowerStr(n int) (string, error) {
	return randomStr(pool[10:36], n)

}

// NoLowerStr create a random string has not lower
func NoLowerStr(n int) (string, error) {
	s := make([]byte, 36)
	copy(s, pool[:10])
	copy(s[10:], pool[36:])
	return randomStr(s, n)
}

// NoNumberStr create a random string has not number
func NoNumberStr(n int) (string, error) {
	return randomStr(pool[10:], n)
}

// randomStr new a random str
func randomStr(source []byte, n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	l := len(source)
	for k, v := range b {
		b[k] = source[int(v)%l]
	}
	return string(b), nil
}
