package binding

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"fmt"

	"github.com/chinx/mohist/web"
	"gopkg.in/yaml.v2"
)

// A request body as multipart/form-data is parsed and up to a total of maxMemory bytes of
// its file parts are stored in memory, with the remainder stored on
// disk in temporary files.
var (
	MaxMemory  int64 = 10 << 20 // 10MB. Should probably make this configurable...
	bindingErr       = "Classification %s error: %s, field name is %s"
)

type (
	Handler interface{}
	parser  func(reflect.Value, *http.Request) error
)

func bind(h Handler, rw http.ResponseWriter, req *http.Request, param web.Params) {
	parse, err := chooseBinder(req)
	if err != nil {

	}
	typ := reflect.TypeOf(h)
	in := make([]reflect.Value, typ.NumIn())
	for i := 0; i < typ.NumIn(); i++ {
		argType := typ.In(i)
		switch argType.String() {
		case "http.ResponseWriter":
			in[i] = reflect.ValueOf(rw)
		case "*http.Request":
			in[i] = reflect.ValueOf(req)
		case "web.Params":
			in[i] = reflect.ValueOf(param)
		default:
			if argType.Kind() != reflect.Ptr {
				panic("Pointers are accepted as bingding models")
			}
			userStruct := reflect.New(argType.Elem())
			parse(userStruct, req, Errors{})
			fromUrl(userStruct, param, Errors{})
			in[i] = userStruct
		}
	}
	callback(reflect.ValueOf(h).Call(in), rw)
}

func errorHandler(errs Errors, rw http.ResponseWriter) {
	if len(errs) > 0 {
		rw.Header().Set("Content-Type", _JSON_CONTENT_TYPE)
		if errs.Has(ERR_DESERIALIZATION) {
			rw.WriteHeader(http.StatusBadRequest)
		} else if errs.Has(ERR_CONTENT_TYPE) {
			rw.WriteHeader(http.StatusUnsupportedMediaType)
		} else {
			rw.WriteHeader(STATUS_UNPROCESSABLE_ENTITY)
		}
		errOutput, _ := json.Marshal(errs)
		rw.Write(errOutput)
		return
	}
}

func chooseBinder(req *http.Request) (b parser, err error) {
	b = fromForm
	contentType := req.Header.Get("Content-Type")
	if req.Method == "POST" || req.Method == "PUT" || len(contentType) > 0 {
		switch {
		case strings.Contains(contentType, "form-urlencoded"):
			b = fromForm
		case strings.Contains(contentType, "multipart/form-data"):
			b = fromMultipart
		case strings.Contains(contentType, "json"):
			b = fromJson
		case strings.Contains(contentType, "yaml"):
			b = fromYaml
		default:
			if contentType == "" {
				err = fmt.Errorf(bindingErr, ERR_CONTENT_TYPE, "Empty Content-Type", "")
			} else {
				err = fmt.Errorf(bindingErr, ERR_CONTENT_TYPE, "Unsupported Content-Type", "")
			}
		}
	}
	return b
}

func Bind(h Handler) web.Handle {
	ensureMethod(h)
	checkHandlerType(reflect.TypeOf(h))
	return func(rw http.ResponseWriter, req *http.Request, param web.Params) {
		bind(h, rw, req, param)
	}
}

func checkHandlerType(typ reflect.Type) {
	for i := 0; i < typ.NumIn(); i++ {
		argType := typ.In(i)
		switch argType.String() {
		case "http.ResponseWriter":
		case "*http.Request":
		case "web.Params":
		default:
			if argType.Kind() != reflect.Ptr {
				panic("Pointers are accepted as bingding models")
			}
		}
	}
}

func fromForm(formStruct reflect.Value, req *http.Request, errors Errors) {
	parseErr := req.ParseForm()

	// Format validation of the request body or the URL would add considerable overhead,
	// and ParseForm does not complain when URL encoding is off.
	// Because an empty request body or url can also mean absence of all needed values,
	// it is not in all cases a bad request, so let's return 422.
	if parseErr != nil {
		errors.Add([]string{}, ERR_DESERIALIZATION, parseErr.Error())
	}
	mapForm(formStruct, req.Form, nil, errors)
}

func fromMultipart(formStruct reflect.Value, req *http.Request, errors Errors) {
	if req.MultipartForm == nil {
		// Workaround for multipart forms returning nil instead of an error
		// when content is not multipart; see https://code.google.com/p/go/issues/detail?id=6334
		if multipartReader, err := req.MultipartReader(); err != nil {
			errors.Add([]string{}, ERR_DESERIALIZATION, err.Error())
		} else {
			form, parseErr := multipartReader.ReadForm(MaxMemory)
			if parseErr != nil {
				errors.Add([]string{}, ERR_DESERIALIZATION, parseErr.Error())
			}

			if req.Form == nil {
				req.ParseForm()
			}
			for k, v := range form.Value {
				req.Form[k] = append(req.Form[k], v...)
			}

			req.MultipartForm = form
		}
	}
	mapForm(formStruct, req.MultipartForm.Value, req.MultipartForm.File, errors)
}

func fromJson(jsonStruct reflect.Value, req *http.Request, errors Errors) {
	if req.Body != nil {
		defer req.Body.Close()
		err := json.NewDecoder(req.Body).Decode(jsonStruct.Interface())
		if err != nil && err != io.EOF {
			errors.Add([]string{}, ERR_DESERIALIZATION, err.Error())
		}
	}
}

func fromYaml(yamlStruct reflect.Value, req *http.Request, errors Errors) {
	if req.Body != nil {
		defer req.Body.Close()
		byts, err := ioutil.ReadAll(req.Body)
		if err == nil {
			err = yaml.Unmarshal(byts, yamlStruct.Interface())
		}
		if err != nil && err != io.EOF {
			errors.Add([]string{}, ERR_DESERIALIZATION, err.Error())
		}
	}
}

func fromUrl(paramsStruct reflect.Value, params web.Params, errors Errors) {
	if paramsStruct.Kind() == reflect.Ptr {
		paramsStruct = paramsStruct.Elem()
	}
	typ := paramsStruct.Type()

	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := paramsStruct.Field(i)

		if typeField.Type.Kind() == reflect.Ptr && typeField.Anonymous {
			structField.Set(reflect.New(typeField.Type.Elem()))
			fromUrl(structField.Elem(), params, errors)
			if reflect.DeepEqual(structField.Elem().Interface(), reflect.Zero(structField.Elem().Type()).Interface()) {
				structField.Set(reflect.Zero(structField.Type()))
			}
		} else if typeField.Type.Kind() == reflect.Struct {
			fromUrl(structField, params, errors)
		}

		inputFieldName := parseFormName(typeField.Name, typeField.Tag.Get("url"))
		if len(inputFieldName) == 0 || !structField.CanSet() {
			continue
		}

		inputValue, exists := params.Get(inputFieldName)
		if exists {
			setWithProperType(typeField.Type.Kind(), inputValue, structField, inputFieldName, errors)
		}
	}
}

func callback(vals []reflect.Value, rw http.ResponseWriter) {
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

// NameMapper represents a form tag name mapper.
type NameMapper func(string) string

var (
	nameMapper = func(field string) string {
		newstr := make([]rune, 0, len(field))
		for i, chr := range field {
			if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
				if i > 0 {
					newstr = append(newstr, '_')
				}
				chr -= ('A' - 'a')
			}
			newstr = append(newstr, chr)
		}
		return string(newstr)
	}
)

func parseFormName(raw, actual string) string {
	if len(actual) > 0 {
		return actual
	}
	return nameMapper(raw)
}

// Takes values from the form data and puts them into a struct
func mapForm(formStruct reflect.Value, form map[string][]string, formfile map[string][]*multipart.FileHeader, errors Errors) {

	if formStruct.Kind() == reflect.Ptr {
		formStruct = formStruct.Elem()
	}
	typ := formStruct.Type()

	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := formStruct.Field(i)

		if typeField.Type.Kind() == reflect.Ptr && typeField.Anonymous {
			structField.Set(reflect.New(typeField.Type.Elem()))
			mapForm(structField.Elem(), form, formfile, errors)
			if reflect.DeepEqual(structField.Elem().Interface(), reflect.Zero(structField.Elem().Type()).Interface()) {
				structField.Set(reflect.Zero(structField.Type()))
			}
		} else if typeField.Type.Kind() == reflect.Struct {
			mapForm(structField, form, formfile, errors)
		}

		inputFieldName := parseFormName(typeField.Name, typeField.Tag.Get("form"))
		if len(inputFieldName) == 0 || !structField.CanSet() {
			continue
		}

		inputValue, exists := form[inputFieldName]
		if exists {
			numElems := len(inputValue)
			if structField.Kind() == reflect.Slice && numElems > 0 {
				sliceOf := structField.Type().Elem().Kind()
				slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
				for i := 0; i < numElems; i++ {
					setWithProperType(sliceOf, inputValue[i], slice.Index(i), inputFieldName, errors)
				}
				formStruct.Field(i).Set(slice)
			} else {
				setWithProperType(typeField.Type.Kind(), inputValue[0], structField, inputFieldName, errors)
			}
			continue
		}

		inputFile, exists := formfile[inputFieldName]
		if !exists {
			continue
		}
		fhType := reflect.TypeOf((*multipart.FileHeader)(nil))
		numElems := len(inputFile)
		if structField.Kind() == reflect.Slice && numElems > 0 && structField.Type().Elem() == fhType {
			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
			for i := 0; i < numElems; i++ {
				slice.Index(i).Set(reflect.ValueOf(inputFile[i]))
			}
			structField.Set(slice)
		} else if structField.Type() == fhType {
			structField.Set(reflect.ValueOf(inputFile[0]))
		}
	}
}

// This sets the value in a struct of an indeterminate type to the
// matching value from the request (via Form middleware) in the
// same type, so that not all deserialized values have to be strings.
// Supported types are string, int, float, and bool.
func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value, nameInTag string, errors Errors) {
	switch valueKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val == "" {
			val = "0"
		}
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			errors.Add([]string{nameInTag}, ERR_INTERGER_TYPE, "Value could not be parsed as integer")
		} else {
			structField.SetInt(intVal)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val == "" {
			val = "0"
		}
		uintVal, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			errors.Add([]string{nameInTag}, ERR_INTERGER_TYPE, "Value could not be parsed as unsigned integer")
		} else {
			structField.SetUint(uintVal)
		}
	case reflect.Bool:
		if val == "on" {
			structField.SetBool(true)
			return
		}

		if val == "" {
			val = "false"
		}
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			errors.Add([]string{nameInTag}, ERR_BOOLEAN_TYPE, "Value could not be parsed as boolean")
		} else if boolVal {
			structField.SetBool(true)
		}
	case reflect.Float32:
		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			errors.Add([]string{nameInTag}, ERR_FLOAT_TYPE, "Value could not be parsed as 32-bit float")
		} else {
			structField.SetFloat(floatVal)
		}
	case reflect.Float64:
		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			errors.Add([]string{nameInTag}, ERR_FLOAT_TYPE, "Value could not be parsed as 64-bit float")
		} else {
			structField.SetFloat(floatVal)
		}
	case reflect.String:
		structField.SetString(val)
	}
}
