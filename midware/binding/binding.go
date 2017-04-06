package binding

import (
	"reflect"

	"net/http"

	"log"
	"strings"

	"github.com/chinx/mohist/web"
)

type Handler interface{}
type binder func(*http.Request, web.Params, Handler)

func bind(rw http.ResponseWriter, req *http.Request, param web.Params, h Handler) {
	typ := reflect.TypeOf(h)
	in := make([]reflect.Value, typ.NumIn())
	for i := 0; i < typ.NumIn(); i++ {
		argType := typ.In(i)
		log.Println(argType)
		switch argType.String() {
		case "http.ResponseWriter":
			in[i] = reflect.ValueOf(rw)
		case "*http.Request":
			in[i] = reflect.ValueOf(req)
		case "web.Params":
			in[i] = reflect.ValueOf(param)
		default:
			userStruct := reflect.New(argType)
			if userStruct.Kind() == reflect.Ptr {
				userStruct = userStruct.Elem()
			}
			chooseBinder(req)(userStruct, req, param)
			in[i] = userStruct
		}
	}
	InvokeBack(reflect.ValueOf(h).Call(in), rw)
}

func chooseBinder(req *http.Request) binder {
	b := Form
	contentType := req.Header.Get("Content-Type")
	if req.Method == "POST" || req.Method == "PUT" || len(contentType) > 0 {
		switch {
		case strings.Contains(contentType, "form-urlencoded"):
			b = Form
		case strings.Contains(contentType, "multipart/form-data"):
			b = MultipartForm
		case strings.Contains(contentType, "json"):
			b = Json
		case strings.Contains(contentType, "yaml"):
			b = Yaml
		default:
		}
	}
	return b
}

func Bind(h Handler) web.Handle {
	ensureMethod(h)
	return func(rw http.ResponseWriter, req *http.Request, param web.Params) {
		bind(rw, req, param, h)
	}
}

func Form(userStruct interface{}, req *http.Request, params web.Params) {

}

func MultipartForm(userStruct interface{}, req *http.Request, params web.Params) {

}

func Json(userStruct interface{}, req *http.Request, params web.Params) {

}

func Yaml(userStruct interface{}, req *http.Request, params web.Params) {

}

func InvokeBack(vals []reflect.Value, rw http.ResponseWriter) {
	// if the handler returned something, write it to the http response
	var respVal reflect.Value
	if len(vals) > 1 && vals[0].Kind() == reflect.Int {
		rw.WriteHeader(int(vals[0].Int()))
		respVal = vals[1]
	} else if len(vals) > 0 {
		respVal = vals[0]

		if isError(respVal) {
			err := respVal.Interface().(error)
			if err != nil {
				rw.Write([]byte(respVal.String()))
			}
			return
		} else if canDeref(respVal) {
			if respVal.IsNil() {
				return // Ignore nil error
			}
		}
	}
	if canDeref(respVal) {
		respVal = respVal.Elem()
	}
	if isByteSlice(respVal) {
		rw.Write(respVal.Bytes())
	} else {
		rw.Write([]byte(respVal.String()))
	}
}

func ensureNotPointer(obj interface{}) {
	if reflect.ValueOf(obj).Kind() == reflect.Ptr {
		panic("Pointers are not accepted as bingding models")
	}
}

func ensureMethod(h Handler) {
	if reflect.TypeOf(h).Kind() != reflect.Func {
		panic("Binding handler must be a callable function")
	}
}

func canDeref(val reflect.Value) bool {
	return val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr
}

func isError(val reflect.Value) bool {
	_, ok := val.Interface().(error)
	return ok
}

func isByteSlice(val reflect.Value) bool {
	return val.Kind() == reflect.Slice && val.Type().Elem().Kind() == reflect.Uint8
}
