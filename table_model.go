package main

type TableModel struct {
	BaseModel
	fields_count int
}

func (table TableModel) Content(int) string {
	return table.name
}

func (table TableModel) EntryPoint() *Window {
	var window *Window
	return window
}
