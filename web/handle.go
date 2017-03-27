package web

import (
	"net/http"
	"net/url"
)

type Handle func(rw http.ResponseWriter, req *http.Request, params *url.Values)

func chainHandler(handlers ...Handle) Handle {
	return func(rw http.ResponseWriter, req *http.Request, params *url.Values) {
		length := len(handlers)
		for i := 0; i < length; i++ {
			if _, ok := rw.(ResponseWriter); ok && rw.(ResponseWriter).Written() {
				break
			}
			handlers[i](rw, req, params)
		}
	}
}
