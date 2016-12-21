package models

import (
	"database/sql"
	adpt "github.com/wanzysky/isqi/adapters"
	"strconv"
)

type TableModel struct {
	BaseModel
	fields []FieldModel
}

func (table TableModel) Content(int) string {
	return table.Name
}

func (table *TableModel) SyncColumns() error {
	if len(table.fields) > 0 {
		return nil
	}

	adapter := adpt.Adpt
	query_sql := adapter.ShowColumns(table.Name, true)
	accepter := make([]interface{}, FIELD_MODEL_ATTR_COUNT)
	containers := make([]sql.RawBytes, FIELD_MODEL_ATTR_COUNT)
	for i := 0; i < FIELD_MODEL_ATTR_COUNT; i++ {
		accepter[i] = &containers[i]
	}

	err := adapter.Query(query_sql, func() {
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
	}, accepter...)
	return err
}

func (table *TableModel) Structure() [][]string {
	table.SyncColumns()
	results := make([][]string, 0)
	header := []string{"field", "type", "collation", "null", "key", "defaults", "extra", "privileges", "comment"}
	results = append(results, header)
	for _, field := range table.fields {
		result := make([]string, 0)
		result = append(result, field.field)
		result = append(result, field.types)
		result = append(result, field.collation)
		result = append(result, field.null)
		result = append(result, field.key)
		result = append(result, field.defaults)
		result = append(result, field.extra)
		result = append(result, field.privileges)
		result = append(result, field.comment)
		results = append(results, result)
	}
	return results
}

func (table *TableModel) Glimpse() [][]string {
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

	adapter.Query(query, func() {
		row := make([]string, count)
		for i, _ := range accepter {
			row[i] = string(containers[i])
		}
		result = append(result, row)
	}, accepter...)
	return result
}

func (table *TableModel) Statistic() (count int) {
	var reciever sql.RawBytes
	var err error
	err = adpt.Adpt.Query(adpt.Adpt.Count(table.Name), func() {
		count, err = strconv.Atoi(string(reciever))
	}, &reciever)
	if err != nil {
		panic(err.Error())
	}
	return
}
