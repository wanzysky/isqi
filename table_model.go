package main

type TableModel struct {
	BaseModel
	fields_count int
}

func (table TableModel) Content(int) string {
	return table.name
}

func (table TableModel) EntryPoint() *Window {
	view := NewTableShowView(table)
	return view.window
}

func (table TableModel) Glimpse() [][]string {
	sql := adapter.Select(table.name)
	result, err := table.Query(sql)
	if err != nil {
		panic("Faild to glimpse table")
	} else {
		return result
	}
}

func (table TableModel) Query(sql string) (rows [][]string, err error) {

	return [][]string{}, nil
}
