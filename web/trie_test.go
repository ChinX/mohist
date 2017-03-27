package web

import (
	"net/http"
	"net/url"
	"testing"
)

func TestAddMatch(t *testing.T) {
	n := newNode()
	n.add("/accounts/", func(w http.ResponseWriter, req *http.Request, params *url.Values) {
		t.Log("one")
	})
	n.add("accounts/:account", func(w http.ResponseWriter, req *http.Request, params *url.Values) {
		t.Log("two")
	})
	n.add("accounts/:account/projects", func(w http.ResponseWriter, req *http.Request, params *url.Values) {
		t.Log("three")
	})
	n.add("accounts/:account/projects/:project", func(w http.ResponseWriter, req *http.Request, params *url.Values) {
		t.Log("four")
	})

	n.add("accounts/:account/projects/abcde", func(w http.ResponseWriter, req *http.Request, params *url.Values) {
		t.Log("five")
	})

	n.add("accounts/:account/projects/:project/files/?:file", func(w http.ResponseWriter, req *http.Request, params *url.Values) {
		t.Log("six")
	})

	param := &url.Values{}
	if handler := n.match("/accounts/", param); handler != nil {
		handler(nil, nil, nil)
	}
	t.Log(param)

	param = &url.Values{}
	if handler := n.match("accounts/account", param); handler != nil {
		handler(nil, nil, nil)
	}
	t.Log(param)

	param = &url.Values{}
	if handler := n.match("accounts/abc/projects", param); handler != nil {
		handler(nil, nil, nil)
	}
	t.Log(param)

	param = &url.Values{}
	if handler := n.match("accounts/account/projects/project", param); handler != nil {
		handler(nil, nil, nil)
	}
	t.Log(param)

	param = &url.Values{}
	if handler := n.match("accounts/account/projects/abcde", param); handler != nil {
		handler(nil, nil, nil)
	}
	t.Log(param)

	param = &url.Values{}
	if handler := n.match("acc/account/projects/abcde", param); handler != nil {
		handler(nil, nil, nil)
	}
	t.Log(param)

	param = &url.Values{}
	if handler := n.match("accounts/account/projects/project/files", param); handler != nil {
		handler(nil, nil, nil)
	}
	t.Log(param)

	param = &url.Values{}
	if handler := n.match("accounts/account/projects/project/files/file", param); handler != nil {
		handler(nil, nil, nil)
	}
	t.Log(param)

	param = &url.Values{}
	if handler := n.match("accounts/account/projects/project/files/file/1", param); handler != nil {
		handler(nil, nil, nil)
	}
	t.Log(param)

}

func matchHandler(w http.ResponseWriter, req *http.Request, param *url.Values) {

}
