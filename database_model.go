package main

import "fmt"

type DatabaseModel struct {
	BaseModel
	tables_count int
	tables       []*TableModel
}

func (db DatabaseModel) Content(length int) string {
	var str string
	if length <= 0 {
		return db.name
	}
	content_format := fmt.Sprintf("%%-%ds", length)
	str = fmt.Sprintf(content_format, db.name)
	return str
}

func (db DatabaseModel) EntryPoint() *Window {
	connection = adapter.Use(db)
	db.tables = Tables()
	tables := NewTableIndexView(db, nav.Position(), nav.Size())
	return tables.window
}
