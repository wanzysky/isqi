package main

import (
	"database/sql"

	ui "github.com/wanzysky/termui"
)

var nav *Navigation
var connection *sql.DB
var adapter Adapter

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	adapter.username = "root"
	adapter.passwd = "root"
	adapter.host = "localhost"
	adapter.port = "3306"

	connection = adapter.Connection()
	databases := Databases(connection)
	defer connection.Close()

	database_view_list := []ItemView{}
	for _, db := range databases {
		database_view_list = append(database_view_list, ItemView{object: db})
	}

	width := ui.TermWidth()
	height := ui.TermHeight()

	main_view := NewListView(0, 3, width, height-3, "Select DataBase", database_view_list)

	dash := NewDashboardView(0, 0, width, 3)
	dash.delegate = main_view
	operatios := map[string]string{
		"s":     "Search",
		"c":     "Quick Choose",
		"d":     "Database Detail",
		"C-c":   "Quit",
		"Enter": "Use",
	}

	tips_str := ""
	for key, op := range operatios {
		tips_str += "[" + key + "] " + "[" + op + "]" + "(fg-white,bg-blue)  "
	}
	dash.tips = tips_str
	window := NewWindow(dash, main_view)
	nav = NewNavigatoin(window)
	nav.Push(window)

	ui.Loop()
}
