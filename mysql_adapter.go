package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Adapter struct {
	username string
	passwd   string
	host     string
	port     string
}

func (adapter Adapter) Connection() *sql.DB {
	dsn_describer := fmt.Sprintf("%s:%s@tcp(%s:%s)/", adapter.username, adapter.passwd, adapter.host, adapter.port)
	connection, err := sql.Open("mysql", dsn_describer)
	if err != nil {
		fmt.Println("Failed to connect Database")
		panic(err.Error())
	}

	return connection
}

func Databases(connection *sql.DB) []DatabaseModel {
	rows, err := connection.Query("SHOW DATABASES")
	if err != nil {
		fmt.Println("Can't connect to host")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	values := make([]sql.RawBytes, 1)
	scan_arg := make([]interface{}, 1)
	var databases []DatabaseModel
	scan_arg[0] = &values[0]

	for rows.Next() {
		err = rows.Scan(scan_arg...)
		if err != nil {
			panic(err.Error())
		}
		value := string(values[0])
		database := DatabaseModel{}
		database.name = value
		databases = append(databases, database)
	}
	return databases
}

func (adapter Adapter) Use(db DatabaseModel) *sql.DB {
	dsn_describer := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", adapter.username, adapter.passwd, adapter.host, adapter.port, db.name)
	connection, err := sql.Open("mysql", dsn_describer)
	if err != nil {
		fmt.Println("Failed to choose Database")
		panic(err.Error())
	}
	return connection
}

func Tables() []*TableModel {
	rows, err := connection.Query("SHOW TABLES")
	var tables []*TableModel
	if err != nil {
		fmt.Println("Can't show tables")
		fmt.Println(err.Error())
		os.Exit(1)
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
		table := TableModel{}
		table.name = value
		tables = append(tables, &table)
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
