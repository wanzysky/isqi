package adapters

import (
	"database/sql"
	"fmt"
	ui "github.com/gizak/termui"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

type MysqlAdapter struct {
	Username string
	Passwd   string
	Host     string
	Port     string
	Conn     *sql.DB
}

func (adapter *MysqlAdapter) Initialize(params map[string]string) {
	adapter.Username = params["username"]
	adapter.Passwd = params["passwd"]
	adapter.Port = params["port"]
	adapter.Host = params["host"]
}

func (adapter *MysqlAdapter) Connect() {
	adapter.Use("")
}

func (adapter *MysqlAdapter) Close() {
	adapter.Conn.Close()
}

func (adapter *MysqlAdapter) Databases() []string {
	defer func() {
		if r := recover(); r != nil {
			ui.Close()
			fmt.Println(r)
			os.Exit(0)
		}
	}()

	rows, err := adapter.Conn.Query("SHOW DATABASES")
	if err != nil {
		log.Panic(err.Error())
	}

	values := make([]sql.RawBytes, 1)
	scan_arg := make([]interface{}, 1)
	var databases []string
	scan_arg[0] = &values[0]

	for rows.Next() {
		err = rows.Scan(scan_arg...)
		if err != nil {
			panic(err.Error())
		}
		value := string(values[0])
		databases = append(databases, value)
	}
	return databases
}

func (adapter *MysqlAdapter) Use(name string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("!")
			ui.Close()
		}
	}()

	dsn_describer := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", adapter.Username, adapter.Passwd, adapter.Host, adapter.Port, name)
	connection, err := sql.Open("mysql", dsn_describer)
	if err != nil {
		panic(err.Error())
	}
	adapter.Conn = connection
}

func (adapter *MysqlAdapter) Tables() []string {
	rows, err := adapter.Conn.Query("SHOW TABLES")
	var tables []string
	if err != nil {
		log.Fatal(err.Error())
		panic("Can't show tables")
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

func (adapter *MysqlAdapter) FullColumns(table string) ([]string, []map[string]string) {
	rows, err := adapter.Conn.Query(ShowColumns(table, true))
	if err != nil {
		log.Panic(err.Error())
	}
	defer rows.Close()
	container := make([]sql.RawBytes, 9)
	accepter := make([]interface{}, 9)
	for i, _ := range container {
		accepter[i] = &(container[i])
	}
	var names []string
	var attrs []map[string]string
	for rows.Next() {
		err = rows.Scan(accepter...)
		if err != nil {
			log.Panic(err.Error())
		}
		name := string(container[0])
		attr := make(map[string]string)
		attr["data type"] = string(container[1])
		attr["collection"] = string(container[2])
		attr["null"] = string(container[3])
		attr["key"] = string(container[4])
		attr["default"] = string(container[5])
		attr["extra"] = string(container[6])
		attr["privileges"] = string(container[7])
		attr["comment"] = string(container[8])
		names = append(names, name)
		attrs = append(attrs, attr)
	}
	return names, attrs
}

func (adapter *MysqlAdapter) Execute(sql string) error {
	_, err := adapter.Conn.Query(sql)
	//var result [][]string
	return err
}

func (adapter *MysqlAdapter) Select(query string) ([][]string, error) {
	rows, err := adapter.Conn.Query(query)
	if err != nil {
		return [][]string{}, err
	}
	defer rows.Close()
	columns, e := rows.Columns()
	if e != nil {
		return [][]string{}, err
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
