package main

type TableIndexView struct {
	BaseView
	tables []*TableModel
	window *Window
}

func NewTableIndexView(db DatabaseModel, position Point, size Size) *TableIndexView {
	table_list := TableIndexView{tables: db.tables}
	var tables []ItemView
	for _, t := range db.tables {
		tables = append(tables, ItemView{object: t})
	}
	list := NewListView(position.x, position.y+3, size.width, size.height-3, db.name, tables)
	dash := NewDashboardView(position.x, position.y, size.width, 3)
	operatios := map[string]string{
		"s":     "Search",
		"c":     "Quick Choose",
		"d":     "Table Detail",
		"C-c":   "Quit",
		"Enter": "Use",
	}

	tips_str := ""
	for key, op := range operatios {
		tips_str += "[" + key + "] " + "[" + op + "]" + "(fg-white,bg-blue)  "
	}
	dash.tips = tips_str
	dash.delegate = list
	table_list.window = NewWindow(dash, list)
	return &table_list
}
