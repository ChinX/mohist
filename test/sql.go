package ormosia

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

type Model struct {
	db        *sql.DB
	tableName string
	param     []string
	columns   string
	where     string
	pk        string
	orderBy   string
	limit     string
	join      string
}

//use goini read the configuration file and connect the mysql database
func SetConfig(username,password,hostname,port,database,charset string) (*Model, error) {
	c := new(Model)
	db, err := sql.Open("mysql", username+":"+password+"@tcp("+hostname+":"+port+")/"+database+"?charset="+charset)
	err = db.Ping()
	if err != nil {
		//if connect error then return the error message
		return c, err
	}
	c.db = db
	return c, err
}

func (m *Model) findSQL() string {
	if len(m.param) == 0 {
		m.columns = "*"
	} else {
		if len(m.param) == 1 {
			m.columns = m.param[0]
		} else {
			m.columns = strings.Join(m.param, ",")
		}

	}
	return fmt.Sprintf("Select %v from %v %v %v %v %v", m.columns, m.tableName, m.join, m.where, m.orderBy, m.limit)
}

func (m *Model) FindAll(result interface{}) error {
	if m.db == nil {
		return fmt.Errorf("mysql not connect")
	}
	rows, err := m.db.Query(m.findSQL())
	if err != nil {
		return err
	}
	return Rows2Structs(rows, result)
}

func (m *Model) FindOne(result interface{}) error {
	if m.db == nil {
		return fmt.Errorf("mysql not connect")
	}
	m.Limit(1)
	rows, err := m.db.Query(m.findSQL())
	if err != nil {
		return err
	}
	return Rows2Struct(rows, result)
}

func (m *Model) Insert(param map[string]interface{}) (num int, err error) {
	if m.db == nil {
		return 0, errors.New("mysql not connect")
	}
	var keys []string
	var values []string
	if len(m.pk) != 0 {
		delete(param, m.pk)
	}
	for key, value := range param {
		keys = append(keys, key)
		switch value.(type) {
		case int, int64, int32:
			values = append(values, strconv.Itoa(value.(int)))
		case string:
			values = append(values, value.(string))
		case float32, float64:
			values = append(values, strconv.FormatFloat(value.(float64), 'f', -1, 64))

		}

	}
	fileValue := "'" + strings.Join(values, "','") + "'"
	fileds := "`" + strings.Join(keys, "`,`") + "`"
	sql := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", m.tableName, fileds, fileValue)
	result, err := m.db.Exec(sql)
	if err != nil {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("SQL syntax errors ")
			}
		}()
		err = errors.New("inster sql failure")
		return 0, err
	}
	i, err := result.LastInsertId()
	s, _ := strconv.Atoi(strconv.FormatInt(i, 10))
	if err != nil {
		err = errors.New("insert failure")
	}
	return s, err

}

func (m *Model) Fields(param ...string) *Model {
	m.param = param
	return m
}

func (m *Model) Update(param map[string]interface{}) (num int, err error) {
	if m.db == nil {
		return 0, errors.New("mysql not connect")
	}
	var setValue []string
	for key, value := range param {
		switch value.(type) {
		case int, int64, int32:
			set := fmt.Sprintf("%v = %v", key, value.(int))
			setValue = append(setValue, set)
		case string:
			set := fmt.Sprintf("%v = '%v'", key, value.(string))
			setValue = append(setValue, set)
		case float32, float64:
			set := fmt.Sprintf("%v = '%v'", key, strconv.FormatFloat(value.(float64), 'f', -1, 64))
			setValue = append(setValue, set)
		}

	}
	setData := strings.Join(setValue, ",")
	sql := fmt.Sprintf("UPDATE %v SET %v %v", m.tableName, setData, m.where)
	result, err := m.db.Exec(sql)
	if err != nil {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("SQL syntax errors ")
			}
		}()
		err = errors.New("update sql failure")
		return 0, err
	}
	i, err := result.RowsAffected()
	if err != nil {
		err = errors.New("update failure")
		return 0, err
	}
	s, _ := strconv.Atoi(strconv.FormatInt(i, 10))

	return s, err
}

//func (m *Model) Delete(param string) (num int, err error) {
//	if m.db == nil {
//		return 0, errors.New("mysql not connect")
//	}
//	h := m.Where(param).FindOne()
//	if len(h) == 0 {
//		return 0, errors.New("no Value")
//	}
//	sql := fmt.Sprintf("DELETE FROM %v WHERE %v", m.tableName, param)
//	result, err := m.db.Exec(sql)
//	if err != nil {
//		defer func() {
//			if err := recover(); err != nil {
//				fmt.Printf("SQL syntax errors ")
//			}
//		}()
//		err = errors.New("delete sql failure")
//		return 0, err
//	}
//	i, err := result.RowsAffected()
//	s, _ := strconv.Atoi(strconv.FormatInt(i, 10))
//	if i == 0 {
//		err = errors.New("delete failure")
//	}
//
//	return s, err
//}

func (m *Model) Count() (int64, error) {
	sql := fmt.Sprintf("Select COUNT(*) from %v %v %v %v %v", m.tableName, m.join, m.where, m.orderBy, m.limit)
	result, err := m.db.Exec(sql)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//func (m *Model) Query(sql string, result interface{}) error {
//	if m.db == nil {
//		return errors.New("mysql not connect")
//	}
//	var query = strings.TrimSpace(sql)
//	s, err := regexp.MatchString(`(?i)^select`, query)
//	if err == nil && s == true {
//		rows, _ := m.db.Query(sql)
//		return QueryResult(rows, result)
//	}
//	exec, err := regexp.MatchString(`(?i)^(update|delete)`, query)
//	if err == nil && exec == true {
//		m_exec, err := m.db.Exec(query)
//		if err != nil {
//			return err
//		}
//		num, _ := m_exec.RowsAffected()
//		id := strconv.FormatInt(num, 10)
//		return id
//	}
//
//	insert, err := regexp.MatchString(`(?i)^insert`, query)
//	if err == nil && insert == true {
//		m_exec, err := m.db.Exec(query)
//		if err != nil {
//			return err
//		}
//		num, _ := m_exec.LastInsertId()
//		id := strconv.FormatInt(num, 10)
//		return id
//	}
//	result, _ := m.db.Exec(query)
//
//	return result
//
//}

func (m *Model) SetTable(tablename string) *Model {
	m.tableName = tablename
	return m
}

func (m *Model) Where(param string) *Model {
	m.where = fmt.Sprintf(" where %v", param)
	return m
}

func (m *Model) SetPk(pk string) *Model {
	m.pk = pk
	return m
}

func (m *Model) OrderBy(param string) *Model {
	m.orderBy = fmt.Sprintf("ORDER BY %v", param)
	return m
}

func (m *Model) Limit(size ...int) *Model {
	var end int
	start := size[0]
	if len(size) > 1 {
		end = size[1]
		m.limit = fmt.Sprintf("Limit %d,%d", start, end)
		return m
	}
	m.limit = fmt.Sprintf("Limit %d", start)
	return m
}

func (m *Model) LeftJoin(table, condition string) *Model {
	m.join = fmt.Sprintf("LEFT JOIN %v ON %v", table, condition)
	return m
}

func (m *Model) RightJoin(table, condition string) *Model {
	m.join = fmt.Sprintf("RIGHT JOIN %v ON %v", table, condition)
	return m
}

func (m *Model) Join(table, condition string) *Model {
	m.join = fmt.Sprintf("INNER JOIN %v ON %v", table, condition)
	return m
}

func (m *Model) FullJoin(table, condition string) *Model {
	m.join = fmt.Sprintf("FULL JOIN %v ON %v", table, condition)
	return m
}

//the function will use friendly way to print the data
func Print(slice map[int]map[string]string) {
	for _, v := range slice {
		for key, value := range v {
			fmt.Println(key, value)
		}
		fmt.Println("---------------")
	}
}

func (m *Model) DbClose() {
	m.db.Close()
}
