package models

import (
	"fmt"
	adpt "github.com/wanzysky/isqi/adapters"
)

type DatabaseModel struct {
	BaseModel
	tables_count int
	Tables       []*TableModel
}

func Databases() (databases []DatabaseModel) {
	names := adpt.Adpt.Databases()
	for _, name := range names {
		db := DatabaseModel{}
		db.Name = name
		databases = append(databases, db)
	}
	return
}

func (db DatabaseModel) Content(length int) string {
	var str string
	if length <= 0 {
		return db.Name
	}
	content_format := fmt.Sprintf("%%-%ds", length)
	str = fmt.Sprintf(content_format, db.Name)
	return str
}

func (db *DatabaseModel) Use() {
	adpt.Adpt.Use(db.Name)
}

func (db *DatabaseModel) FetchTables() []*TableModel {
	for _, table_name := range adpt.Adpt.Tables() {
		var table_model TableModel
		table_model.Name = table_name
		db.Tables = append(db.Tables, &table_model)
	}
	return db.Tables
}
