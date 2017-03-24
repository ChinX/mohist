package main

import (
	"log"
	"net/http"

	"github.com/chinx/mohist/web"
	"github.com/urfave/negroni"
)

func main() {
	r := web.NewRouter(web.NotFound)
	InitRouter(r)

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":9090")
}

func InitRouter(r *web.Router) {
	r.Get("/accounts/", testHandle)
	r.Get("/accounts/:account", testHandle)
	r.Get("/accounts/:account/status", testHandle)
	r.Get("/accounts/:account/status/:statu", testHandle)
}

func testHandle(w http.ResponseWriter, req *http.Request) {
	log.Println(req.URL, req.Form)
	w.Write([]byte("ok"))
	w.WriteHeader(200)
}
