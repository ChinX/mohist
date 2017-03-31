package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chinx/mohist/web"
	"github.com/urfave/negroni"
	"net/url"
)

func main() {
	r := web.NewRouter()
	InitRouter(r)

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":9999")
}

func InitRouter(r *web.Router) {
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
				r.Get("/", firstHandler, testHandle)
				r.Get("/:statu", firstHandler, testHandle)
			}, groupThreeHandler)
		}, groupTwoHandler)
	}, groupFirstHandler)
}

func groupFirstHandler(w http.ResponseWriter, req *http.Request, params *url.Values) {
	log.Println("this is first grou")
}

func groupTwoHandler(w http.ResponseWriter, req *http.Request, params *url.Values) {
	log.Println("this is two grou")
}

func groupThreeHandler(w http.ResponseWriter, req *http.Request, params *url.Values) {
	log.Println("this is three grou")
}

func firstHandler(w http.ResponseWriter, req *http.Request, params *url.Values) {
	log.Println("this is frist")
}

func testHandle(w http.ResponseWriter, req *http.Request, params *url.Values) {
	backStr := fmt.Sprintf("%s: %s", req.URL.Path, params)
	log.Println(backStr)
	w.WriteHeader(200)
	w.Write([]byte(backStr))
}
