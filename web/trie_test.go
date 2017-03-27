package web

import (
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/bmizerany/pat"
	"github.com/gin-gonic/gin"
	"github.com/go-baa/baa"
	"github.com/go-macaron/macaron"
	"github.com/julienschmidt/httprouter"
	"github.com/labstack/echo"
)

type testPat struct {
}

func (p testPat) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

}

func BenchmarkMohist(b *testing.B) {
	n := NewRouter()
	execute(b, n, func(path string) {
		n.Get(path, func(w http.ResponseWriter, req *http.Request, params *url.Values) {
		})
	})
}

func BenchmarkEach(b *testing.B) {
	n := echo.New()
	execute(b, n, func(path string) {
		n.GET(path, func(content *echo.Context) {
		})
	})
}

func BenchmarkGin(b *testing.B) {
	n := gin.New()
	execute(b, n, func(path string) {
		n.GET(path, func(content *gin.Context) {
		})
	})
}

func BenchmarkBaa(b *testing.B) {
	n := baa.New()
	execute(b, n, func(path string) {
		n.Get(path, func(content *baa.Context) {
		})
	})
}

func BenchmarkPat(b *testing.B) {
	n := pat.New()
	execute(b, n, func(path string) {
		n.Get(path, &testPat{})
	})
}

func BenchmarkHttpRouter(b *testing.B) {
	n := &httprouter.Router{}
	execute(b, n, func(path string) {
		n.GET(path, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {})
	})
}

func BenchmarkMacaron(b *testing.B) {
	n := macaron.New()
	execute(b, n, func(path string) {
		n.Get(path, func(w http.ResponseWriter, req *http.Request) {})
	})
}

func execute(b *testing.B, n http.Handler, handler func(string)) {
	addUrls := addUrl()
	for i := 0; i < len(addUrls); i++ {
		handler(addUrls[i])
	}
	matches := matchUrl()
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(matches); i++ {
			request(n, matches[i])
		}
	}
}

func request(n http.Handler, path string) {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		log.Println(err)
	}
	n.ServeHTTP(nil, req)
}

func addUrl() []string {
	return []string{
		"/accounts/",
		"/accounts/:account",
		"/accounts/:account/projects",
		"/accounts/:account/projects/:project",
		"/accounts/:account/projects/:project/files/:file",
	}
}

func matchUrl() []string {
	return []string{
		"/accounts/",
		"/accounts/account",
		"/accounts/account/projects",
		"/accounts/account/projects/project",
		"/accounts/account/projects/project/files/file",
	}
}
