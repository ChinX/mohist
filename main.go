package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chinx/mohist/web"
	"github.com/urfave/negroni"
)

func main() {
	r := web.NewRouter()
	initRouter(r)
	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":9999")
}

func initRouter(r *web.Router) {
	r.Group("/accounts", func() {
		r.Get("/", firstHandler, testHandle)
		r.Group("/:account", func() {
			r.Get("/", firstHandler, testHandle)
			r.Group("/status", func() {
				r.Get("/", firstHandler, testHandle)
				r.Get("/:statu", firstHandler, testHandle)
			}, groupThreeHandler)
		}, groupTwoHandler)
	}, groupFirstHandler)
	r.Group("/ccounts", func() {
		r.Get("/", firstHandler, testHandle)
		r.Group("/:account", func() {
			r.Get("/", firstHandler, testHandle)
			r.Group("/status", func() {
				r.Get("/aaaa", firstHandler, testaaaaHandle)
				r.Get("/:bbb", firstHandler, testbbbHandle)
				r.Get("/:ccc/abc", firstHandler, testcccHandle)
				r.Get("/?:statu", firstHandler, testdddHandle)
			}, groupThreeHandler)
		}, groupTwoHandler)
	}, groupFirstHandler)
}

func groupFirstHandler(w http.ResponseWriter, req *http.Request, params web.Params) {
	log.Println("this is first grou")
}

func groupTwoHandler(w http.ResponseWriter, req *http.Request, params web.Params) {
	log.Println("this is two grou")
}

func groupThreeHandler(w http.ResponseWriter, req *http.Request, params web.Params) {
	log.Println("this is three grou")
}

func firstHandler(w http.ResponseWriter, req *http.Request, params web.Params) {
	log.Println("this is frist")
}

func testHandle(w http.ResponseWriter, req *http.Request, params web.Params) {
	backStr := fmt.Sprintf("%s: %s", req.URL.Path, params)
	log.Println(backStr)
	w.WriteHeader(200)
	w.Write([]byte(backStr))
}

func testaaaaHandle(w http.ResponseWriter, req *http.Request, params web.Params) {
	backStr := fmt.Sprintf("%s: %s params %s", req.URL.Path, "aaaa", params)
	log.Println(backStr)
	w.WriteHeader(200)
	w.Write([]byte(backStr))
}

func testbbbHandle(w http.ResponseWriter, req *http.Request, params web.Params) {
	backStr := fmt.Sprintf("%s: %s params %s", req.URL.Path, "bbb", params)
	log.Println(backStr)
	w.WriteHeader(200)
	w.Write([]byte(backStr))
}

func testcccHandle(w http.ResponseWriter, req *http.Request, params web.Params) {
	backStr := fmt.Sprintf("%s: %s params %s", req.URL.Path, "ccc", params)
	log.Println(backStr)
	w.WriteHeader(200)
	w.Write([]byte(backStr))
}

func testdddHandle(w http.ResponseWriter, req *http.Request, params web.Params) {
	backStr := fmt.Sprintf("%s: %s params %s", req.URL.Path, "?", params)
	log.Println(backStr)
	w.WriteHeader(200)
	w.Write([]byte(backStr))
}
