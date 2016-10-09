package main

type DatabaseModel struct {
	BaseModel
	tables_count int
	tables       [100]*TableModel
}
