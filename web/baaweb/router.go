package baaweb

import (
	"net/http"
	"net/url"
)

const (
	GET int = iota
	POST
	PUT
	DELETE
	PATCH
	OPTIONS
	HEAD
	// RouteLength route table length
	RouteLength
)

// RouterMethods declare method key in routeMap
var RouterMethods = map[string]int{
	http.MethodGet:     GET,
	http.MethodPost:    POST,
	http.MethodPut:     PUT,
	http.MethodDelete:  DELETE,
	http.MethodPatch:   PATCH,
	http.MethodOptions: OPTIONS,
	http.MethodHead:    HEAD,
}

// A request body as multipart/form-data is parsed and up to a total of maxMemory bytes of
// its file parts are stored in memory, with the remainder stored on
// disk in temporary files.
var maxMemory int64 = 10 << 20 // 10MB. Should probably make this configurable...

type (
	group struct {
		pattern  string
		handlers []Handle
	}
)

type Params []struct {
	Key   string
	Value string
}

func (ps Params) Get(key string) (string, bool) {
	for _, entry := range ps {
		if entry.Key == key {
			return entry.Value, true
		}
	}
	return "", false
}

type Router struct {
	Trees   map[string]*Node
	notFond Handle
	groups  []*group
}

func NewRouter() *Router {
	return &Router{
		Trees: make(map[string]*Node, RouteLength),
	}
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(maxMemory)

	path := req.URL.Path
	nrw := NewResponseWriter(rw)
	params := &url.Values{}
	if root, ok := r.Trees[req.Method]; ok {
		if handler := root.match(path, params); handler != nil {
			handler(nrw, req, params)
			return
		}
	}
	// Handle 404
	if r.notFond != nil {
		r.notFond(rw, req, nil)
		return
	}
	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte(path + " not fond"))
}

func (r *Router) NotFound(handlers ...Handle) {
	r.notFond = handlersChain(handlers)
}

func (r *Router) Group(pattern string, fn func(), handlers ...Handle) {
	r.groups = append(r.groups, &group{"/" + TrimByte(pattern, '/'), handlers})
	fn()
	r.groups = r.groups[:len(r.groups)-1]
}

func (r *Router) Head(pattern string, handlers ...Handle) {
	r.handle(http.MethodHead, pattern, handlers)
}

func (r *Router) Get(pattern string, handlers ...Handle) {
	r.handle(http.MethodGet, pattern, handlers)
	r.handle(http.MethodHead, pattern, handlers)
}

func (r *Router) Post(pattern string, handlers ...Handle) {
	r.handle(http.MethodPost, pattern, handlers)
}

func (r *Router) Put(pattern string, handlers ...Handle) {
	r.handle(http.MethodPut, pattern, handlers)
}

func (r *Router) Delete(pattern string, handlers ...Handle) {
	r.handle(http.MethodDelete, pattern, handlers)
}

func (r *Router) Patch(pattern string, handlers ...Handle) {
	r.handle(http.MethodPatch, pattern, handlers)
}

func (r *Router) Options(pattern string, handlers ...Handle) {
	r.handle(http.MethodOptions, pattern, handlers)
}

func (r *Router) Any(pattern string, handlers ...Handle) {
	r.handle(http.MethodGet, pattern, handlers)
	r.handle(http.MethodPost, pattern, handlers)
	r.handle(http.MethodHead, pattern, handlers)
	r.handle(http.MethodDelete, pattern, handlers)
	r.handle(http.MethodPatch, pattern, handlers)
	r.handle(http.MethodOptions, pattern, handlers)
}

func (r *Router) handle(method, pattern string, handlers []Handle) {
	pattern = "/" + TrimByte(pattern, '/')
	if len(r.groups) > 0 {
		groupPattern := ""
		h := make([]Handle, 0)
		for _, g := range r.groups {
			groupPattern += g.pattern
			h = append(h, g.handlers...)
		}
		pattern = groupPattern + pattern
		h = append(h, handlers...)
		handlers = h
	}
	root, ok := r.Trees[method]
	if !ok {
		root = newNode()
		r.Trees[method] = root
	}
	root.addNode(pattern, handlersChain(handlers))
}

func handlersChain(handlers []Handle) Handle {
	nHandlers := make([]Handle, 0, 10)
	l := 0
	for i := 0; i < len(handlers); i++ {
		if handlers[i] != nil {
			nHandlers = append(nHandlers, handlers[i])
			l++
		}
	}
	if l == 0 {
		return nil
	}
	return func(rw http.ResponseWriter, req *http.Request, params *url.Values) {
		length := len(handlers)
		for i := 0; i < length; i++ {
			if _, ok := rw.(ResponseWriter); ok && rw.(ResponseWriter).Written() {
				return
			}
			handlers[i](rw, req, params)
		}
		if _, ok := rw.(ResponseWriter); ok && !rw.(ResponseWriter).Written() {
			rw.Write([]byte("Mohist is OK"))
		}
	}
}
