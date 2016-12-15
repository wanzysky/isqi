package models

import (
	"database/sql"
	adpt "isqi/adapters"
)

type TableModel struct {
	BaseModel
	fields []FieldModel
}

func (table TableModel) Content(int) string {
	return table.Name
}

func (table *TableModel) SyncColumns() error {
	adapter := adpt.Adpt
	query_sql := adapter.ShowColumns(table.Name, true)
	accepter := make([]interface{}, FIELD_MODEL_ATTR_COUNT)
	containers := make([]sql.RawBytes, FIELD_MODEL_ATTR_COUNT)
	for i := 0; i < FIELD_MODEL_ATTR_COUNT; i++ {
		accepter[i] = &containers[i]
	}

	err := adapter.Query(query_sql, accepter, func(acpt []interface{}) {
		field := FieldModel{}
		field.field = string(containers[0])
		field.types = string(containers[1])
		field.collation = string(containers[2])
		field.null = string(containers[3])
		field.key = string(containers[4])
		field.defaults = string(containers[5])
		field.extra = string(containers[6])
		field.privileges = string(containers[7])
		field.comment = string(containers[8])
		table.fields = append(table.fields, field)
	})
	return err
}

//func (table TableModel) EntryPoint() *Window {
//	view := NewTableShowView(table)
//	return view.window
//}

func (table TableModel) Glimpse() [][]string {
	adapter := adpt.Adpt
	table.SyncColumns()
	query := adapter.Select(table.Name)
	result := [][]string{}
	header := []string{}
	count := len(table.fields)
	for _, field := range table.fields {
		header = append(header, field.field)
	}
	result = append(result, header)

	accepter := make([]interface{}, count)
	containers := make([]sql.RawBytes, count)
	for i, _ := range accepter {
		accepter[i] = &containers[i]
	}

	adapter.Query(query, accepter, func(acpt []interface{}) {
		row := make([]string, count)
		for i, _ := range acpt {
			row[i] = string(containers[i])
		}
		result = append(result, row)
	})
	return result
}
