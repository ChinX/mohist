package rflct

import (
	"fmt"
	"reflect"
	"strconv"
)

func ValOfNotNilPtr(i interface{}) reflect.Value {
	rv := reflect.ValueOf(i)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		rt := reflect.TypeOf(i)
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

func IndirectType(t reflect.Type) reflect.Type {
	if t.Kind() != reflect.Ptr {
		return t
	}
	return t.Elem()
}

func IndirectTypeVal(rv reflect.Value) (reflect.Type, reflect.Value) {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return rv.Type(), rv
}

func GetValByTagName(tag, name string, rv reflect.Value) string {
	rt := rv.Type()
	rt.Key()
	for i := 0; i < rt.NumField(); i++ {
		rsi := rt.Field(i)
		if rsi.Anonymous || rsi.Tag.Get(tag) != name {
			continue
		}
		rvi := rv.Field(i)
		if rvi.Kind() != reflect.Ptr || !rvi.IsNil() && rvi.CanSet() {
			str := ValString(rvi)

			return str
		}
	}
	return ""
}

func SetValTagPair(rv reflect.Value, tag string, keys, vals []string) error {
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
						if err := SetValParseString(rvi, vals[ki]); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func SetValParseString(v reflect.Value, s string) (err error) {
	rt, v := IndirectTypeVal(v)
	switch rt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var result int64
		result, err = strconv.ParseInt(s, 10, 64)
		if err == nil {
			v.SetInt(result)
		}
	case reflect.Bool:
		var result bool
		result, err = strconv.ParseBool(s)
		if err == nil {
			v.SetBool(result)
		}
	case reflect.Float32, reflect.Float64:
		var result float64
		result, err = strconv.ParseFloat(s, 0)
		if err == nil {
			v.SetFloat(result)
		}
	case reflect.String:
		v.SetString(s)
	default:
		return fmt.Errorf("file struct error(not found value)")
	}
	return nil
}

func ValString(v reflect.Value) string {
	rt, v := IndirectTypeVal(v)

	switch rt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', 2, 64)
	case reflect.String:
		return v.String()
	}
	return ""
}

func FieldsIndex(t reflect.Type, tagName string) map[string][]int {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fields := make(map[string][]int)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get(tagName)

		// Skip unexported fields or fields marked with "-"
		if f.PkgPath != "" || tag == "-" {
			continue
		}
		// Handle embedded structs
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			for k, v := range FieldsIndex(f.Type, tagName) {
				fields[k] = append(f.Index, v...)
			}
			continue
		}

		if tag == "" {
			tag = snakeCasedName(f.Name)
		}

		fields[tag] = f.Index
	}
	return fields
}
