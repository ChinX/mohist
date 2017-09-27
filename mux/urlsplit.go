package mux

import "fmt"

var strictErr = "Url \"%s\" has repeated '/' on strict mode "

// UrlSplit if strict is false, url "/domains//domain///users////user"
// will be think of it as "/domains/domain/users/user"
func UrlSplit(url string, strict bool) (parts []string, err error) {
	parts = []string{}
	start, length := 1, len(url)
	switch length {
	case start:
	case start + 1:
		if url[start] == '/' {
			if strict {
				err = fmt.Errorf(strictErr, url)
			}
		} else {
			parts = append(parts, url[start:])
		}
	default:
		index, ended:=0, length
		prefix, suffix := false, false
		for index <= ended {
			if !suffix {
				last := ended - 1
				if last < index{
					break
				}
				suffix = bool(url[last] != '/')
				if !suffix{
					ended = last
					continue
				}
				if strict && length - ended > 2{
					err = fmt.Errorf(strictErr, url)
					break
				}
			}

			delimiter := bool(url[index] == '/')
			if !prefix && !delimiter {
				start = index
				prefix = true
			} else if prefix && delimiter {
				parts = append(parts, url[start:index])
				prefix = false
			}
			index = index + 1
		}
		if start != ended {
			parts = append(parts, url[start:ended])
		}

		//path := url + "a"
		//index, border := 1, false
		//for ; index < length+1; index++ {
		//	delimiter := bool(path[index] == '/')
		//	if delimiter && border {
		//		if strict {
		//			err = fmt.Errorf(strictErr, url)
		//			start = length
		//			break
		//		}
		//	} else if delimiter != border {
		//		if delimiter && start != index {
		//			parts = append(parts, path[start:index])
		//		}
		//		start = index
		//		border = delimiter
		//	}
		//}
		//if start != length {
		//	parts = append(parts, path[start:length])
		//}
	}
	return
}

type UrlParse struct {
	url    string
	index  int
	length int
}

func NewParse(url string) *UrlParse {
	return &UrlParse{url: url, index: 0, length: len(url)}
}

func (u *UrlParse) Traverse() (part string, ending bool) {
	index := u.index
	switch u.length {
	case index + 1:
		if u.url[index] != '/' {
			part = u.url
		}
		fallthrough
	case index:
		index, ending = u.length, true
	default:
		start, begin := index, false
		for ; index < u.length; index++ {
			if (u.url[index] == '/') == begin {
				if begin {
					break
				} else {
					start = index
					begin = true
				}
			}
		}
		ending = bool(u.length == index)
		if index-1 > start {
			part = u.url[start:index]
		}
	}
	u.index = index
	return
}
