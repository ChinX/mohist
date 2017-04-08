package binding

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/chinx/mohist/web"
)

type Sname struct {
	A   string `form:"a"`
	B   string `form:"b"`
	AAA string `url:"aaa"`
	CCC string `url:"ccc"`
}

func TestBind(t *testing.T) {
	handler := Bind(func(obj Sname) (int, []byte) {
		log.Println(obj)
		return 0, []byte("aaa")
	})
	rw := httptest.NewRecorder()
	s := &Sname{A: "name", B: "value"}
	byt, _ := json.Marshal(s)
	req := httptest.NewRequest("GET", "http://127.0.0.1/ass/abc?a=thisisa&b=thisisb", bytes.NewReader(byt))
	params := web.Params{&web.param{Key: "aaa", Value: "bbb"}, &web.param{Key: "ccc", Value: "ddd"}}
	handler(rw, req, params)
}
