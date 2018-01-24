package rflct

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ValPointerNotNil(val interface{}) reflect.Value {
	rv := reflect.ValueOf(val)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		rt := reflect.TypeOf(val)
		msg := "reflect mode is nil"
		if rt != nil {
			if rt.Kind() != reflect.Ptr {
				msg = "reflect mode (" + rt.String() + ") is non-pointer"
			} else {
				msg = "reflect mode (" + rt.String() + ") is nil"
			}
		}
		panic(msg)
	}
	return rv.Elem()
}

func realTypeValue(rv reflect.Value) (reflect.Type, reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return rv.Type(), rv
}

func valByTagName(tag, name string, rv reflect.Value) string {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		rsi := rt.Field(i)
		if rsi.Anonymous || rsi.Tag.Get(tag) != name {
			continue
		}
		rvi := rv.Field(i)
		if rvi.Kind() != reflect.Ptr || !rvi.IsNil() && rvi.CanSet() {
			str := strWithProperType(rvi)

			return str
		}
	}
	return ""
}

func setByTagPair(rv reflect.Value, tag string, keys, vals []string) error {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		rsi := rt.Field(i)
		if rsi.Anonymous {
			continue
		}

		if tagName := rsi.Tag.Get(tag); tagName != "" {
			for ki := range keys {
				if keys[ki] == tagName {
					rvi := rv.Field(i)
					if rvi.Kind() != reflect.Ptr || !rvi.IsNil() && rvi.CanSet() {
						if err := setWithProperType(rvi, vals[ki]); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func setWithProperType(rv reflect.Value, val string) (err error) {
	rt, rv := realTypeValue(rv)
	switch rt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var result int64
		result, err = strconv.ParseInt(val, 10, 64)
		if err == nil {
			rv.SetInt(result)
		}
	case reflect.Bool:
		var result bool
		result, err = strconv.ParseBool(val)
		if err == nil {
			rv.SetBool(result)
		}
	case reflect.Float32, reflect.Float64:
		var result float64
		result, err = strconv.ParseFloat(val, 0)
		if err == nil {
			rv.SetFloat(result)
		}
	case reflect.String:
		rv.SetString(val)
	default:
		return fmt.Errorf("file struct error(not found value)")
	}
	return nil
}

func strWithProperType(rv reflect.Value) string {
	rt, rv := realTypeValue(rv)

	switch rt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'f', 2, 64)
	case reflect.String:
		return rv.String()
	}
	return ""
}

func FieldInfo(typ reflect.Type, tagName string) map[string][]int {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	fields := make(map[string][]int)
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		tag := f.Tag.Get(tagName)

		// Skip unexported fields or fields marked with "-"
		if f.PkgPath != "" || tag == "-" {
			continue
		}
		// Handle embedded structs
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			for k, v := range FieldInfo(f.Type, tagName) {
				fields[k] = append(f.Index, v...)
			}
			continue
		}

		// Use field name for untagged fields
		if tag == "" {
			tag = f.Name
		}

		tag = strings.ToLower(tag)

		fields[tag] = f.Index
	}
	return fields
}
