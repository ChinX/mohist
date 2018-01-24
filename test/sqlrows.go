package ormosia

import (
	"github.com/chinx/ormosia/rflct"
	"reflect"
	"strings"
	"database/sql"
)

func RowsUnmarshStruct(rows *sql.Rows, result interface{}) error  {
	columns, err := rows.Columns()
	if err == nil {
		values := make([]string, len(columns))
		scans := make([]interface{}, len(columns))
		for i := range values {
			scans[i] = &values[i]
		}

		isScan := false
		for rows.Next() {
			if isScan {
				continue
			}
			rows.Scan(scans...)
			rflct.Unmarshal(result, "orm", columns, values)
			isScan = true
		}
	}
	rows.Close()
	return rows.Err()
}

func RowsUnmarshStructs(rows *sql.Rows, result interface{}) error  {
	columns, err := rows.Columns()
	if err == nil {
		values := make([]string, len(columns))
		scans := make([]interface{}, len(columns))
		for i := range values {
			scans[i] = &values[i]
		}

		for rows.Next() {
			rows.Scan(scans...)
			rflct.Unmarshal(result, "orm", columns, values)
		}
	}
	rows.Close()
	return rows.Err()
}

func Rows2Struct(rows *sql.Rows, result interface{}) error  {
	columns, err := rows.Columns()
	if err == nil {
		v := rflct.ValPointerNotNil(result)

		t := v.Type()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		fieldsInfo := rflct.FieldInfo(t, "orm")

		var values []interface{}
		for _, column := range columns {
			idx, ok := fieldsInfo[strings.ToLower(column)]
			var val interface{}
			if !ok {
				var i interface{}
				val = &i
			} else {
				val = v.FieldByIndex(idx).Addr().Interface()
			}
			values = append(values, val)
		}
		isScan := false
		for rows.Next() {
			if isScan {
				continue
			}
			rows.Scan(values...)
			isScan = true
		}
	}
	rows.Close()
	return rows.Err()
}

func Rows2Structs(rows *sql.Rows, result interface{}) error  {
	columns, err := rows.Columns()
	if err == nil {
		v := rflct.ValPointerNotNil(result)
		t := v.Type().Elem()

		fieldsInfo := rflct.FieldInfo(t, "orm")

		for rows.Next() {
			var rv reflect.Value
			var fv reflect.Value

			if t.Kind() == reflect.Ptr {
				rv = reflect.New(t.Elem())
				fv = reflect.Indirect(rv)
			} else {
				rv = reflect.Indirect(reflect.New(t))
				fv = rv
			}

			var values []interface{}
			for _, column := range columns {
				idx, ok := fieldsInfo[strings.ToLower(column)]
				var val interface{}
				if !ok {
					var i interface{}
					val = &i
				} else {
					val = fv.FieldByIndex(idx).Addr().Interface()
				}
				values = append(values, val)
			}
			rows.Scan(values...)
			rflct.SliceExpandSet(v, rv)
		}
	}
	rows.Close()
	return rows.Err()
}
