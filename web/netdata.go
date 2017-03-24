package web

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondBytes(rw http.ResponseWriter, val []byte, status int) {
	rw.Write(val)
	rw.WriteHeader(status)
}

func RespondString(rw http.ResponseWriter, val string, status int) {
	rw.Header().Set("Content-Type", "text/html;charset=utf-8")
	RespondBytes(rw, []byte(val), status)
}

func RespondJson(rw http.ResponseWriter, val interface{}, status int) {
	byteArr, err := json.Marshal(val)
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to marshal to json: %s", err), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json;charset=utf-8")
	RespondBytes(rw, byteArr, status)
}

func SetJsonHeader(header http.Header, key string, val interface{}) error {
	byt, err := json.Marshal(val)
	if err != nil {
		return err
	}
	header.Set(key, base64.URLEncoding.EncodeToString(byt))
	return nil
}

func SetBase64Header(header http.Header, key string, val string) error {
	header.Set(key, base64.URLEncoding.EncodeToString([]byte(val)))
	return nil
}

func getBase64Header(header http.Header, key string) (string, error) {
	return base64.URLEncoding.DecodeString(header.Get(key))
}

func GetBase64Header(req *http.Request, key string) (string, error) {
	byt, err := getBase64Header(req, key)
	return string(byt), err
}

func GetJsonHeader(req *http.Request, key string, val interface{}) error {
	byt, err := getBase64Header(req, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(byt, val)
}
