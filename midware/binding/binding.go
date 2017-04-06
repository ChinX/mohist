package binding

import "reflect"

type bindHandle func(obj interface{}) (int, []byte)

func Bind(obj interface{}, binder bindHandle) /*web.Handle*/ {
	ensureNotPointer(obj)
	nStruct := reflect.New(reflect.TypeOf(obj))
	if nStruct.Kind() == reflect.Ptr{
		nStruct = nStruct.Elem()
	}
	typ := nStruct.Type()
	for i := 0; i < typ.NumField(); i++{
		typField := typ.Field(i)
		sField := nStruct.Field(i)
		if typField.Type.Kind() == reflect.Ptr && typField.Anonymous{
			sField.Set(reflect.New(typField.Type.Elem()))
		}else if typField.Type.Kind() == reflect.String{
			sField.SetString("abc")
		}
	}
	//binder(obj)
	binder(nStruct.Interface())
}

func ensureNotPointer(obj interface{}) {
	if reflect.ValueOf(obj).Kind() == reflect.Ptr {
		panic("Pointers are not accepted as bingding models")
	}
}
