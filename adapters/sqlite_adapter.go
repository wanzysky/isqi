package adapters

import (
	"database/sql"
	"fmt"
	ui "github.com/gizak/termui"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type SqliteAdapter struct {
	path string
	Conn *sql.DB
}

func (adapter *SqliteAdapter) Initialize(params map[string]string) {
	adapter.path = params["file"]
}

func (adapter *SqliteAdapter) Connect() {
	adapter.Use("")
}

func (adapter *SqliteAdapter) Close() {
	adapter.Conn.Close()
}

func (adapter *SqliteAdapter) Databases() []string {
	return make([]string, 1)
}

func (adapter *SqliteAdapter) Use(name string) {
	defer func() {
		if r := recover(); r != nil {
			ui.Close()
		}
	}()
	connection, err := sql.Open("sqlite3", adapter.path)
	adapter.Conn = connection
	if err != nil {
		log.Panic(err.Error())
	}
}

func (adapter *SqliteAdapter) Tables() []string {
	rows, err := adapter.Conn.Query("SELECT name FROM sqlite_master WHERE type='table';")
	var tables []string
	if err != nil {
		log.Panic(err.Error())
	}

	values := make([]sql.RawBytes, 1)
	scan_arg := make([]interface{}, 1)
	scan_arg[0] = &values[0]

	for rows.Next() {
		err = rows.Scan(scan_arg...)
		if err != nil {
			panic(err.Error())
		}
		value := string(values[0])
		tables = append(tables, value)
	}
	return tables
}

func (adapter *SqliteAdapter) FullColumns(table string) ([]string, []map[string]string) {
	rows, err := adapter.Conn.Query(fmt.Sprintf("PRAGMA table_info(%s);", table))
	if err != nil {
		log.Panic(err.Error())
	}
	defer rows.Close()
	var columns []string
	var attrs []map[string]string
	for rows.Next() {
		var name string
		var data_type string
		var null string
		var default_value string
		attr := make(map[string]string)
		attr["data type"] = data_type
		attr["null"] = null
		attr["default value"] = default_value

		rows.Scan(&name, &data_type, &null, &default_value)
		columns = append(columns, name)
		attrs = append(attrs, attr)
	}
	return columns, attrs
}

func (adapter *SqliteAdapter) Execute(sql string) error {
	_, err := adapter.Conn.Query(sql)
	//var result [][]string
	return err
}

func (adapter *SqliteAdapter) Select(query string) ([][]string, error) {
	rows, err := adapter.Conn.Query(query)
	if err != nil {
		return [][]string{}, err
	}
	defer rows.Close()
	columns, e := rows.Columns()
	if e != nil {
		return [][]string{}, e
	}
	width := len(columns)
	var results [][]string
	results = append(results, columns)
	container := make([]sql.RawBytes, width)
	accepter := make([]interface{}, width)
	for i, _ := range container {
		accepter[i] = &(container[i])
	}
	count := 0
	for rows.Next() {
		if count > 1000 {
			break
		}
		err = rows.Scan(accepter...)
		if err != nil {
			return results, err
		}
		result := make([]string, width)
		for index, cell := range container {
			result[index] = string(cell)
		}
		count += 1
		results = append(results, result)
	}
	return results, nil
}
