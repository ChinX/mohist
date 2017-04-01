package web

import (
	"log"
	"net/http"
	"testing"

	"net/http/httptest"
	"time"

	"os"

	"encoding/json"
	"net/url"

	"github.com/bmizerany/pat"
	"github.com/gin-gonic/gin"
	"github.com/go-baa/baa"
	"github.com/go-macaron/macaron"
	"github.com/julienschmidt/httprouter"
	"github.com/labstack/echo"
)

var waitTime time.Duration = 0

type testPat struct {
}

func (p testPat) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	wait()
}

func TestRouter_ServeHTTP(t *testing.T) {
	n := newNode()
	//partFunc("//accounts/:account//", func(part string, ending bool) {
	//	log.Println(part)
	//})
	addUrls := addUrl()
	for i := 0; i < len(addUrls); i++ {
		n.addNode(addUrls[i], func(w http.ResponseWriter, req *http.Request, params *url.Values) {
			wait()
		})
		//n.addNode(addUrls[i], Handle(addUrls[i]))

	}
	matchUrls := matchUrl()
	for i := 0; i < len(matchUrls); i++ {
		handler := n.match(matchUrls[i], &url.Values{})

		log.Println(handler)

	}
	byteArr, err := json.Marshal(n)
	log.Print(string(byteArr), err)
}

func BenchmarkGin(b *testing.B) {
	os.Setenv("GIN_MODE", "1")
	n := gin.New()
	b.Log("Gin")
	execute(b, func(path string) {
		n.GET(path, func(content *gin.Context) {
			wait()
		})
	}, func(path string) { request(n, path) })
}

func BenchmarkMohist(b *testing.B) {
	n := NewRouter()
	b.Log("Mohist")
	execute(b, func(path string) {
		n.Get(path, func(w http.ResponseWriter, req *http.Request, params *url.Values) {
			wait()
		})
		//n.Get(path, Handle(path))
	}, func(path string) { request(n, path) })
}

func BenchmarkEach(b *testing.B) {
	n := echo.New()
	b.Log("Each")
	execute(b, func(path string) {
		n.GET(path, func(content echo.Context) error {
			wait()
			return nil
		})
	}, func(path string) { request(n, path) })
}

func BenchmarkBaa(b *testing.B) {
	n := baa.New()
	b.Log("Baa")
	execute(b, func(path string) {
		n.Get(path, func(content *baa.Context) {
			wait()
		})
	}, func(path string) { request(n, path) })
}

func BenchmarkPat(b *testing.B) {
	n := pat.New()
	b.Log("Pat")
	execute(b, func(path string) {
		n.Get(path, &testPat{})
	}, func(path string) { request(n, path) })
}

func BenchmarkHttpRouter(b *testing.B) {
	n := &httprouter.Router{}
	b.Log("httprouter")
	execute(b, func(path string) {
		n.GET(path, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
			wait()
		})
	}, func(path string) { request(n, path) })
}

func BenchmarkMacaron(b *testing.B) {
	n := macaron.New()
	b.Log("Macaron")
	execute(b, func(path string) {
		n.Get(path, func(w http.ResponseWriter, req *http.Request) {
			wait()
		})
	}, func(path string) { request(n, path) })
}

func execute(b *testing.B, reg func(string), req func(string)) {
	addUrls := addUrl()
	for i := 0; i < len(addUrls); i++ {
		reg(addUrls[i])
	}
	matches := matchUrl()
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(matches); i++ {
			req(matches[i])
		}
	}
}

func request(n http.Handler, path string) {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		log.Println(err)
	}
	n.ServeHTTP(httptest.NewRecorder(), req)
}

func wait() {
	if waitTime == 0 {
		return
	}
	time.Sleep(time.Millisecond * waitTime)
}

func addUrl() []string {
	return []string{
		"/accounts/",
		"/accounts/:account",
		"/accounts/:account/projects",
		"/accounts/:account/projects/:project",
		"/accounts/:account/projects/:project/files/:file",
		"/ccounts/",
		"/ccounts/:account",
		"/ccounts/:account/projects",
		"/ccounts/:account/projects/:project",
		"/ccounts/:account/projects/:project/files/:file",
	}
}

func matchUrl() []string {
	return []string{
		"/accounts/",
		"/accounts/account/",
		"/accounts/account/projects/",
		"/accounts/account/projects/project/",
		"/accounts/account/projects/project/files/file/111",
		"/ccounts/",
		"/ccounts/account/",
		"/ccounts/account/projects/",
		"/ccounts/account/projects/project/",
		"/ccounts/account/projects/project/files/file/",
	}
}
