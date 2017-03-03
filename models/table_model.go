package models

import (
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

	query_sql := adpt.ShowColumns(table.Name, true)

	rows, err := adpt.Adpt.Select(query_sql)
	containers := rows[0]
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
	table.SyncColumns()
	query := adpt.Select(table.Name)

	result, err := adpt.Adpt.Select(query)
	if err != nil {
		panic("Failed to show table")
	}
	return result
}

func (table *TableModel) Statistic() (count int) {
	rows, err := adpt.Adpt.Select(adpt.Count(table.Name))
	if err != nil {
		panic(err)
	}
	count, _ = strconv.Atoi(rows[1][0])
	return count
}
