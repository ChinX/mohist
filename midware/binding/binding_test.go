package binding

import (
	"testing"
	"log"
	"net/http"
	"github.com/chinx/mohist/web"
)

type Sname struct {
	A string
	B string
}

func TestBind(t *testing.T) {
	Bind(func(obj *Sname, rw http.ResponseWriter, req *http.Request, params web.Params)(int, []byte){
		log.Println(obj)
		log.Println(rw)
		log.Println(req)
		log.Println(params)
		return 0, []byte("aaa")
	})
}
