package mux

import (
	"net/http"
	"testing"
	"log"
)

func TestMatch(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	l := &Leaf{
		Handler: func(writer http.ResponseWriter, request *http.Request, params Params) {},
		Children: []*Leaf{
			{
				Pattern: "abc",
				IsParam: false,
				Children: []*Leaf{
					{
						Pattern: "def",
						Handler: func(writer http.ResponseWriter, request *http.Request, params Params) {},
						IsParam: false,
						Children: []*Leaf{
							{
								Pattern: "ghi",
								IsParam: false,
								Children: []*Leaf{
									{
										Pattern: "jkl",
										Handler: func(writer http.ResponseWriter, request *http.Request, params Params) {},
										IsParam: false,
									},
								},
							},
							{
								Pattern: "mno",
								IsParam: true,
								Children: []*Leaf{
									{
										Pattern: "pqr",
										Handler: func(writer http.ResponseWriter, request *http.Request, params Params) {},
										IsParam: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	//t.Log(l.match("") != nil)
	//t.Log(l.match("abc") != nil)
	//t.Log(l.match("abc/def") != nil)
	//t.Log(l.match("abc/def/aaa") != nil)
	t.Log(l.match("abc/def/aaa/bbb") != nil)
	//t.Log(l.match("abc/def/aaa/bbb/ccc") != nil)
}
