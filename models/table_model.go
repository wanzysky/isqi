package models

import (
	adpt "github.com/wanzysky/isqi/adapters"
	"log"
	"strconv"
)

type TableModel struct {
	BaseModel
	fields []FieldModel
}

func (table TableModel) Content(int) string {
	return table.Name
}

func (table *TableModel) SyncColumns() {
	if len(table.fields) > 0 {
		return
	}

	names, attrs := adpt.Adpt.FullColumns(table.Name)
	var fields []FieldModel
	for i, name := range names {
		fields = append(fields, FieldModel{name: name, attributes: attrs[i]})
	}
	table.fields = fields
}

func (table *TableModel) Structure() [][]string {
	table.SyncColumns()
	results := make([][]string, 0)
	if len(table.fields) <= 0 {
		return results
	}
	return results
}

func (table *TableModel) Glimpse() [][]string {
	table.SyncColumns()
	query := adpt.Select(table.Name)

	result, err := adpt.Adpt.Select(query)
	if err != nil {
		log.Panic("Failed to show table")
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
