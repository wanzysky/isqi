package adapters

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Processer func()

type Adapter struct {
	Username string
	Passwd   string
	Host     string
	Port     string
}

var Conn *sql.DB
var Adpt Adapter

func Connection() *sql.DB {
	dsn_describer := fmt.Sprintf("%s:%s@tcp(%s:%s)/", Adpt.Username, Adpt.Passwd, Adpt.Host, Adpt.Port)
	connection, err := sql.Open("mysql", dsn_describer)
	if err != nil {
		fmt.Println("Failed to connect Database")
		panic(err.Error())
	}
	Conn = connection
	return Conn
}

func Databases() []string {
	rows, err := Conn.Query("SHOW DATABASES")
	if err != nil {
		fmt.Println("Can't connect to host")
		fmt.Println(err.Error())
		os.Exit(1)
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

func Use(name string) *sql.DB {
	dsn_describer := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", Adpt.Username, Adpt.Passwd, Adpt.Host, Adpt.Port, name)
	connection, err := sql.Open("mysql", dsn_describer)
	Conn = connection
	if err != nil {
		panic(err.Error())
	}
	return connection
}

func Tables() []string {
	rows, err := Conn.Query("SHOW TABLES")
	var tables []string
	if err != nil {
		fmt.Println("Can't show tables")
		panic(err.Error())
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

func (adapter Adapter) Select(table_name string, args ...string) string {
	var fields string
	if len(args) == 0 {
		fields = "*"
	} else {
		fields = strings.Join(args, ",")
	}

	return "SELECT " + fields + " FROM " + table_name + " LIMIT 100"
}

func (adapter Adapter) Count(table_name string) string {
	return "SELECT COUNT(1) FROM " + table_name
}

func (dtapter Adapter) ShowColumns(table_name string, full bool) string {
	return "SHOW FULL COLUMNS FROM " + table_name
}

func (adapter Adapter) Query(sql string, callback Processer, accepter ...interface{}) error {
	rows, err := Conn.Query(sql)
	//var result [][]string
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(accepter...)
		callback()
	}

	return nil
}
