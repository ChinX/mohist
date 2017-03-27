package web

import (
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/chinx/mohist/web/nrouter"
	"github.com/go-macaron/macaron"
)

func BenchmarkBytes(b *testing.B) {
	//f, err := os.OpenFile("./new/bytes", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	//if err != nil {
	//	b.Fatal(err)
	//}
	//log.SetOutput(f)
	//n := nrouter.NewNode()
	//paths, logs := addUrl()
	//for i := 0; i < len(paths); i++ {
	//	addBytes(n, []byte(paths[i]), logs[i])
	//}

	//for i := 0; i < b.N; i++ {
	//	matches := matchUrl()
	//	for i := 0; i < len(matches); i++ {
	//		matchBytes(n, []byte(paths[i]), &url.Values{})
	//	}
	//}
	//f.Close()

	a := []byte("abc")
	log.Println(string(a[:1]))
	log.Println(string(a))
}

func BenchmarkMatch(b *testing.B) {
	//f, err := os.OpenFile("./new/self", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	//if err != nil {
	//	b.Fatal(err)
	//}
	//log.SetOutput(f)
	n := newNode()
	paths, logs := addUrl()
	for i := 0; i < len(paths); i++ {
		addSelf(n, paths[i], logs[i])
	}

	matches := matchUrl()
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(matches); i++ {
			matchSelf(n, matches[i], &url.Values{})
		}
	}
	//f.Close()
}

func BenchmarkMacaron(b *testing.B) {
	//f, err := os.OpenFile("./new/macaron", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	//if err != nil {
	//	b.Fatal(err)
	//}
	//log.SetOutput(f)
	n := macaron.NewTree()

	addUrls, logs := addUrl()
	for i := 0; i < len(addUrls); i++ {
		addMacaron(n, addUrls[i], logs[i])
	}

	matches := matchUrl()
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(matches); i++ {
			matchMacaron(n, matches[i])
		}
	}
	//f.Close()
}

func addBytes(n *nrouter.Node, path []byte, logs string) {
	n.Add(path, func(w http.ResponseWriter, req *http.Request, params *url.Values) {
		log.Println(logs)
	})
}

func matchBytes(n *nrouter.Node, path []byte, params *url.Values) {
	if handler := n.Match(path, params); handler != nil {
		handler(nil, nil, nil)
		log.Println(params)
	}
}

func addSelf(n *node, path, logs string) {
	n.add(path, func(w http.ResponseWriter, req *http.Request, params *url.Values) {
		//log.Println(logs)
	})
}

func matchSelf(n *node, path string, params *url.Values) {
	if handler := n.match(path, params); handler != nil {
		handler(nil, nil, nil)
		//log.Println(params)
	}
}

func addMacaron(n *macaron.Tree, path, logs string) {
	n.Add(path, func(w http.ResponseWriter, req *http.Request, params macaron.Params) {
		//log.Println(logs)
	})
}

func matchMacaron(n *macaron.Tree, path string) {
	if handler, params, ok := n.Match(path); ok {
		handler(nil, nil, params)
		//log.Println(params)
	}

}

func addUrl() ([]string, []string) {
	return []string{
			"/accounts/",
			"/accounts/:account",
			"/accounts/:account/projects",
			"/accounts/:account/projects/:project",
			"/accounts/:account/projects/abcde",
			"/accounts/:account/projects/:project/files/?:file",
		}, []string{
			"one",
			"two",
			"three",
			"four",
			"five",
			"six",
		}
}

func matchUrl() []string {
	return []string{
		"/accounts/",
		"/accounts/account",
		"/accounts/abc/projects",
		"/accounts/account/projects/project",
		"/accounts/account/projects/abcde",
		"/acc/account/projects/abcde",
		"/acc/account/projects/project/files",
		"/acc/account/projects/project/files/file",
		"/acc/account/projects/project/files/file/1",
	}
}
