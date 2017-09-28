package mux

import (
	"io"
	"log"

	"github.com/chinx/mohist/internal"
)

const (
	static = iota
	dynamic
	elastic
)

type Leaf struct {
	Children  []*Leaf
	Elastomer *Leaf
	Handler   Handle
	Pattern   string
	Kind      int
	IsParam   bool
}

func (l *Leaf) compare(part string) (key string, matched bool) {
	if l.IsParam {
		log.Println("compare, isparam true key:", l.Pattern)
		return l.Pattern, l.IsParam
	}
	log.Println("compare, Pattern == part:", l.Pattern, part)
	return "", l.Pattern == part
}

func (l *Leaf) matchChild(path string, start int) (Handle, error) {
	next, part, err := internal.Traverse(path, start, '/')
	log.Println("child match, next num is:", next, "; this part is :", part, "; the err: ", err)

	if err == nil || err == io.EOF {
		ending := err == io.EOF

		for i := 0; i < len(l.Children); i++ {
			child := l.Children[i]

			key, ok := child.compare(part)
			if ok && ending {
				log.Println("matched params key: ", key, "; it exist", ok, "; Traverse is", ending)
				if ending {
					if key != "" {
						log.Printf("matched leaf param %s => %s", key, part)
					}
					return child.Handler, err
				}
			}

			h, err1 := child.matchChild(path, next)
			if h != nil {
				if key != ""{
					log.Printf("matched leaf param %s => %s", key, part)
				}
			}else if err1 == io.EOF {
				continue
			}
			return h, err1
		}
		if l.Elastomer != nil {
			return l.Elastomer.Handler, nil
		}
	}
	return nil, err
}

func (l *Leaf) match(url string) Handle {
	path := internal.Trim(url, '/')
	if path == "" {
		return l.Handler
	}
	log.Println("in match, this path is", path)
	h, err := l.matchChild(path, 0)
	if err != nil && err != io.EOF {
		return nil
	}
	log.Println("match over, function is", h, "; err is", err)
	return h
}
