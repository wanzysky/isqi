package windows

import (
	"image"
	model "isqi/models"
	view "isqi/views"
)

type TableListWindow struct {
	*ListWindow
	tables []*model.TableModel
}

func NewTableIndexWindow(db *model.DatabaseModel) *TableListWindow {
	var table_list TableListWindow
	var tables []view.ItemView
	db.Use()
	for _, t := range db.FetchTables() {
		tables = append(tables, view.ItemView{Object: t})
	}
	rect := Size()
	list := view.NewListView(image.Rect(rect.Min.X, rect.Min.Y+3, rect.Max.X, rect.Max.Y), db.Name, tables)
	dash := view.NewDashboardView(image.Rect(rect.Min.X, rect.Min.Y, rect.Max.X, 3))
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
	dash.Tips = tips_str
	dash.Delegate = list
	table_list.ListWindow = &ListWindow{Window: NewWindow(dash, list), list: list, dash: dash}
	return &table_list
}
