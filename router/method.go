package router

import (
	"net/http"
	"net/url"
)

func ChainHandler(handlers ...Handle) Handle {
	return func(w http.ResponseWriter, req *http.Request, param *url.Values) {
		length := len(handlers)
		for i := 0; i < length; i++ {
			handlers[i](w, req, param)
		}
	}
}
