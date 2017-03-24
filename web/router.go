package web

import (
	"net/http"
	"strings"
)

var (
	// known HTTP methods.
	httpMethods = map[string]bool{
		http.MethodGet:     true,
		http.MethodPost:    true,
		http.MethodPut:     true,
		http.MethodDelete:  true,
		http.MethodPatch:   true,
		http.MethodOptions: true,
		http.MethodHead:    true,
	}
)

const (
	MethodAny = "*"
)

type Handle func(w http.ResponseWriter, req *http.Request)

type Router struct {
	tree        *node
	rootHandler Handle
}

func NewRouter(rootHandler Handle) *Router {
	node := node{component: "/", isNamedParam: false, methods: make(map[string]*route)}
	return &Router{tree: &node, rootHandler: rootHandler}
}

func (r *Router) eachHandle(method, path string, handler Handle) {
	method = strings.ToUpper(method)
	if !strings.HasPrefix(path, "/") {
		path = "/" + strings.TrimSuffix(path, "/")
	}

	var methods map[string]bool
	if method == MethodAny {
		methods = httpMethods
	} else if httpMethods[method] {
		methods = map[string]bool{method: true}
	} else {
		panic("unknown HTTP method: " + method)
	}

	for m := range methods {
		r.tree.addNode(m, path, handler)
	}
}

func (r *Router) handle(method, pattern string, hs []Handle) {
	for _, h := range hs {
		r.eachHandle(method, pattern, h)
	}
}

func (r *Router) Group(pattern string, fn func(), hs ...Handle) {

}

func (r *Router) Head(pattern string, hs ...Handle) {
	r.handle(http.MethodHead, pattern, hs)
}

func (r *Router) Get(pattern string, hs ...Handle) {
	r.handle(http.MethodGet, pattern, hs)
	//if r.autoHead {
	//	r.Hand(pattern, hs...)
	//}
}

func (r *Router) Patch(pattern string, hs ...Handle) {
	r.handle(http.MethodPatch, pattern, hs)
}

func (r *Router) Post(pattern string, hs ...Handle) {
	r.handle(http.MethodPost, pattern, hs)
}

func (r *Router) Put(pattern string, hs ...Handle) {
	r.handle(http.MethodPut, pattern, hs)
}

func (r *Router) Options(pattern string, hs ...Handle) {
	r.handle(http.MethodOptions, pattern, hs)
}

func (r *Router) Any(pattern string, hs ...Handle) {
	r.handle(MethodAny, pattern, hs)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	cw := w
	req.ParseMultipartForm(10 * 1024 * 1024) // 10MB. Should probably make this configurable...
	params := &req.Form
	//if !runMiddleware(cw, req, params, r.middleware...) {
	//	return // end the chain.
	//}
	node, _ := r.tree.traverse(strings.Split(req.URL.Path, "/")[1:], params)
	if handler := node.methods[req.Method]; handler != nil {
		//if !runMiddleware(cw, req, params, handler.middleware...) {
		//	return
		//}
		handler.handler(cw, req)
	} else {
		r.rootHandler(cw, req)
	}
}

func NotFound(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(req.URL.Path + " not fond"))
	w.WriteHeader(http.StatusNotFound)
}
