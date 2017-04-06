package binding

import (
	"testing"
	"log"
)

type Sname struct {
	A string
	B string
}

func TestBind(t *testing.T) {
	Bind(Sname{}, func(obj interface{})(int, []byte){
		log.Println(obj.(Sname))
		return 0, []byte("aaa")
	})
}
