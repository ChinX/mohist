package web

import (
	"net/http"
	"net/url"
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
	// A request body as multipart/form-data is parsed and up to a total of maxMemory bytes of
	// its file parts are stored in memory, with the remainder stored on
	// disk in temporary files.
	MaxMemory int64 = 10 * 1024 * 1024 // 10MB. Should probably make this configurable...
)

const (
	MethodAny = "*"
)

type group struct {
	pattern  string
	handlers []Handle
}

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

type router struct {
	autoHead       bool
	trees          map[string]*node
	notFondHandler Handle
	groups         []*group
}

func NewRouter() Router {
	return &router{
		trees: make(map[string]*node),
	}
}

func (r *router) handle(method, pattern string, handlers ...Handle) {
	pattern = formatPattern(pattern)
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
	method = strings.ToUpper(method)
	var methods map[string]bool
	if method == MethodAny {
		methods = httpMethods
	} else if httpMethods[method] {
		methods = map[string]bool{method: true}
	} else {
		panic("unknown HTTP method: " + method)
	}

	for m := range methods {
		t, ok := r.trees[m]
		if !ok {
			t = newNode()
			r.trees[m] = t
		}
		t.add(pattern, chainHandler(handlers...))
	}
}

func (r *router) notFound(rw http.ResponseWriter, req *http.Request, params *url.Values) {
	if r.notFondHandler != nil {
		r.notFondHandler(rw, req, params)
		return
	}
	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte(req.URL.Path + " not fond"))
}

func (r *router) NotFound(handlers ...Handle) {
	r.notFondHandler = chainHandler(handlers...)
}

func (r *router) Group(pattern string, fn func(), handlers ...Handle) {
	r.groups = append(r.groups, &group{formatPattern(pattern), handlers})
	fn()
	r.groups = r.groups[:len(r.groups)-1]
}

func (r *router) Head(pattern string, handlers ...Handle) {
	r.handle(http.MethodHead, pattern, handlers...)
}

func (r *router) Get(pattern string, handlers ...Handle) {
	r.handle(http.MethodGet, pattern, handlers...)
	if r.autoHead {
		r.Head(pattern, handlers...)
	}
}

func (r *router) Patch(pattern string, handlers ...Handle) {
	r.handle(http.MethodPatch, pattern, handlers...)
}

func (r *router) Post(pattern string, handlers ...Handle) {
	r.handle(http.MethodPost, pattern, handlers...)
}

func (r *router) Put(pattern string, handlers ...Handle) {
	r.handle(http.MethodPut, pattern, handlers...)
}

func (r *router) Options(pattern string, handlers ...Handle) {
	r.handle(http.MethodOptions, pattern, handlers...)
}

func (r *router) Any(pattern string, handlers ...Handle) {
	r.handle(MethodAny, pattern, handlers...)
}

func (r *router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(MaxMemory)
	params := &url.Values{}
	nrw := NewResponseWriter(rw)
	if t, ok := r.trees[req.Method]; ok {
		if handler := t.match(req.URL.Path, params); handler != nil {
			handler(nrw, req, params)
			return
		}
	}
	r.notFound(nrw, req, params)
}

func formatPattern(pattern string) string {
	return "/" + strings.Trim(pattern, "/")
}
