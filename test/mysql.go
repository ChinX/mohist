package mysql

import (
	"database/sql"
	"fmt"

	"github.com/chinx/dbproxy/svcs"
)

func init() {
	svcs.Registry("mysql", NewService)
}

type Client struct {
	DB        *sql.DB
	tableName string
	where     string
}

func NewService(driver, host, port, user, passwd, dbName string) (svcs.DB, error) {
	hostName := host
	if port != "" {
		hostName += ":" + port
	}

	source := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, passwd, hostName, dbName)
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping db is error")
	}
	return &Client{DB: db}, nil
}

func (c *Client) Table(name string) svcs.Tabulator {
	c.tableName = name
	return c
}

func (c *Client) Where(conditions ...*svcs.Pair) svcs.Behavior {
	//WHERE condition1 [AND [OR]] condition2.....
	if len(conditions) < 1 {
		return c
	}
	where := ""
	for index := range conditions {
		where += " AND " + conditions[index].Key + "=" + conditions[index].Val
	}
	if where != "" {
		c.where = "WHERE " + where[5:]
	}
	return c
}

func (c *Client) Create(values ...*svcs.Pair) (int64, error) {
	//INSERT INTO table_name (column1, column2,...) VALUES (value1, value2,....)
	if len(values) < 1 {
		return 0, fmt.Errorf("create values must not empty")
	}
	keyStr := ""
	valStr := ""
	for index := range values {
		keyStr += "," + values[index].Key
		valStr += "," + values[index].Val
	}
	if keyStr != "" {
		keyStr = keyStr[1:]
		valStr = valStr[1:]
	}
	parse := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", c.tableName, keyStr, valStr)
	res, err := c.DB.Exec(parse)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
func (c *Client) Retrieve(filters ...string) ([][]*svcs.Pair, error) {
	//SELECT field1, field2,...fieldN FROM table_name1, table_name2... [WHERE condition1 [AND [OR]] condition2.....
	filter := ""
	for index := range filters {
		filter += "," + filters[index]
	}

	if filter == "" {
		filter = "*"
	} else {
		filter = filter[1:]
	}

	rowPairs := make([][]*svcs.Pair, 0, 10)
	parse := fmt.Sprintf("SELECT %s FROM %s %s", filter, c.tableName, c.where)
	rows, err := c.DB.Query(parse)
	if err != nil {
		return nil, err
	}
	if columns, err := rows.Columns(); err != nil {
		return nil, err
	} else {
		//拼接记录Map
		values := make([]sql.RawBytes, len(columns))
		scans := make([]interface{}, len(columns))

		for i := range values {
			scans[i] = &values[i]
		}
		//此处遍历在3W记录的时候，长达1分钟甚至更多
		for rows.Next() {
			_ = rows.Scan(scans...)
			each := make([]*svcs.Pair, len(values))
			for i, col := range values {
				each[i] = &svcs.Pair{Key: columns[i], Val: string(col)}
			}
			rowPairs = append(rowPairs, each)
		}
	}
	return rowPairs, nil
}

func (c *Client) Update(values ...*svcs.Pair) (int, error) {
	return 0, nil
}

func (c *Client) Delete() (int, error) {
	return 0, nil
}
