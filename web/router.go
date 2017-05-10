package web

import (
	"net/http"

	"github.com/chinx/mohist/byteconv"
)

type Router interface {
	NotFound(handlers ...Handle)
	Group(pattern string, fn func(), handlers ...Handle)
	Head(pattern string, handlers ...Handle)
	Get(pattern string, handlers ...Handle)
	Patch(pattern string, handlers ...Handle)
	Post(pattern string, handlers ...Handle)
	Put(pattern string, handlers ...Handle)
	Options(pattern string, handlers ...Handle)
	Any(pattern string, handlers ...Handle)
	http.Handler
}

type group struct {
	pattern  string
	handlers []Handle
}

type router struct {
	Trees   map[string]*node
	notFond Handle
	groups  []*group
}

func NewRouter() Router {
	return &router{
		Trees: make(map[string]*node, 7),
	}
}

func (r *router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if root, ok := r.Trees[req.Method]; ok {
		if handler, params := root.match(path); handler != nil {
			handler(NewResponseWriter(rw), req, params)
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
	rw = nil
}

func (r *router) NotFound(handlers ...Handle) {
	r.notFond = handlersChain(handlers)
}

func (r *router) Group(pattern string, fn func(), handlers ...Handle) {
	r.groups = append(r.groups, &group{"/" + byteconv.Trim(pattern, '/'), handlers})
	fn()
	r.groups = r.groups[:len(r.groups)-1]
}

func (r *router) Head(pattern string, handlers ...Handle) {
	r.handle(http.MethodHead, pattern, handlers)
}

func (r *router) Get(pattern string, handlers ...Handle) {
	r.handle(http.MethodGet, pattern, handlers)
	r.handle(http.MethodHead, pattern, handlers)
}

func (r *router) Post(pattern string, handlers ...Handle) {
	r.handle(http.MethodPost, pattern, handlers)
}

func (r *router) Put(pattern string, handlers ...Handle) {
	r.handle(http.MethodPut, pattern, handlers)
}

func (r *router) Delete(pattern string, handlers ...Handle) {
	r.handle(http.MethodDelete, pattern, handlers)
}

func (r *router) Patch(pattern string, handlers ...Handle) {
	r.handle(http.MethodPatch, pattern, handlers)
}

func (r *router) Options(pattern string, handlers ...Handle) {
	r.handle(http.MethodOptions, pattern, handlers)
}

func (r *router) Any(pattern string, handlers ...Handle) {
	r.handle(http.MethodGet, pattern, handlers)
	r.handle(http.MethodPost, pattern, handlers)
	r.handle(http.MethodHead, pattern, handlers)
	r.handle(http.MethodDelete, pattern, handlers)
	r.handle(http.MethodPatch, pattern, handlers)
	r.handle(http.MethodOptions, pattern, handlers)
}

func (r *router) handle(method, pattern string, handlers []Handle) {
	pattern = "/" + byteconv.Trim(pattern, '/')
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
	return func(rw http.ResponseWriter, req *http.Request, params *Params) {
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
