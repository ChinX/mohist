package internal

import (
	"fmt"
	"io"
	"bytes"
)

var strictErr = " '%s' is repeated in source \"%s\""

func UrlSplit(url string) (arr []string, err error) {
	url = Trim(url, '/')
	for start, part := 0, ""; err == nil; {
		start, part, err = Traverse(url, start, '/')
		if err != nil && err != io.EOF {
			return
		}
		arr = append(arr, part)
	}
	err = nil
	return
}

func Traverse(url string, index int, b byte) (next int, part string, err error) {
	length := len(url)
	switch length {
	case index, index + 1:
		part = url[index:]
		next, err = length, io.EOF
	default:
		next = index
		for start := next - 1; next <= length; next++ {
			if next == length || url[next] == b {
				switch next - start {
				case 0:
				case 1:
					start = next
					if start == index || start == length {
						continue
					}
					buff := bytes.NewBuffer([]byte{})
					buff.WriteByte(b)
					err = fmt.Errorf(strictErr, buff.Bytes(), url)
					return
				default:
					part = url[start+1 : next]
					if next == length {
						err = io.EOF
					}
					return
				}
			}
		}
	}
	return
}
