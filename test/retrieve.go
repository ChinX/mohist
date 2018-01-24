func (c *Client) ReadOne(result interface{}, filters ...string) error {
	c.Limit(0, 1)
	rows, err := c.retrieve(filters)
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err == nil {
		values := make([]string, len(columns))
		scans := make([]interface{}, len(columns))
		for i := range values {
			scans[i] = &values[i]
		}

		if rows.Next() && err == nil {
			_ = rows.Scan(scans...)
			err = rflct.Unmarshal(result, "orm", columns, values)
		}
		rows.Close()
	}
	return err
}

func (c *Client) ReadAll(result interface{}, filters ...string) error {
	c.Limit(0, 1)
	rows, err := c.retrieve(filters)
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err == nil {
		values := make([]string, len(columns))
		scans := make([]interface{}, len(columns))

		for i := range values {
			scans[i] = &values[i]
		}

		for rows.Next() && err == nil {
			_ = rows.Scan(scans...)
			err = rflct.Unmarshal(result, "orm", columns, values)
		}
		rows.Close()
	}
	return err
}

func (c *Client) ReadOneNew(result interface{}, filters ...string) error {
	c.Limit(0, 1)
	rows, err := c.retrieve(filters)
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err != nil{
		rows.Close()
		return err
	}

	v := rflct.ValPointerNotNil(result)

	t := v.Type()
	if t.Kind() == reflect.Ptr{
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
	if rows.Next() {
		_ = rows.Scan(values...)
		log.Println(values)
	}
	rows.Close()
	return rows.Err()
}

func (c *Client) ReadAllNew(result interface{}, filters ...string) error {
	rows, err := c.retrieve(filters)
	if err != nil {
		return err
	}
	columns, err := rows.Columns()
	if err != nil{
		rows.Close()
		return err
	}
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
		err = rows.Scan(values...)
		if err != nil {
			return err
		}
		rflct.SliceExpandSet(v, rv)
	}

	rows.Close()
	return rows.Err()
}
