package adapters

import (
	"strings"
)

type Adapter interface {
	Initialize(map[string]string)
	Connect()
	Databases() []string
	Use(string)
	Tables() []string
	Select(string) ([][]string, error)
	Execute(string) error
	Close()
}

var Adpt Adapter

func Select(table_name string, args ...string) string {
	var fields string
	if len(args) == 0 {
		fields = "*"
	} else {
		fields = strings.Join(args, ",")
	}

	return "SELECT " + fields + " FROM " + table_name + " LIMIT 100"
}

func ShowColumns(table_name string, full bool) string {
	return "SHOW FULL COLUMNS FROM " + table_name
}

func Count(table_name string) string {
	return "SELECT COUNT(1) FROM " + table_name
}

func Initialize(adapter_name string) {
	switch adapter_name {
	case "mysql", "mysql2":
		Adpt = &(MysqlAdapter{})
	case "sqlite":
	default:
		panic("Unkown adapter" + adapter_name)
	}
}
