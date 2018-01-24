package rflct

import (
	"bytes"

	"fmt"
	"reflect"

	"github.com/chinx/utils/strutil"
)

func FillRestful(restful string, s interface{}) string {
	rv := ValPointerNotNil(s)
	buffer := &bytes.Buffer{}
	length := len(restful)
	for path, next := "", 0; next < length; {
		path, next = strutil.Between(restful, '/', next)
		if path[0] == ':' {
			path = valByTagName("restful", path[1:], rv)
		}
		buffer.WriteByte('/')
		buffer.WriteString(path)
	}
	return buffer.String()
}

func Unmarshal(s interface{}, tag string, keys, vals []string) error {
	if len(keys) != len(vals) {
		return fmt.Errorf("unmarshal keys and vals must be equal in length")
	}

	v := ValPointerNotNil(s)

	if v.Kind() == reflect.Slice {
		var item reflect.Value
		if t := v.Type(); t.Kind() == reflect.Ptr {
			item = reflect.New(t.Elem())
		} else {
			item = reflect.New(t)
		}
		SliceExpandSet(v, item)
		v = item.Elem()
	}
	return setByTagPair(v, tag, keys, vals)
}

func SliceExpandSet(sv, v reflect.Value) {
	i := sv.Len()
	if i >= sv.Cap() {
		ncap := sv.Cap() + sv.Cap()/2
		if ncap < 4 {
			ncap = 4
		}
		nv := reflect.MakeSlice(sv.Type(), sv.Len(), ncap)
		reflect.Copy(nv, sv)
		sv.Set(nv)
	}
	sv.SetLen(i + 1)
	sv.Index(i).Set(v)
}
