package main

import ui "github.com/gizak/termui"

type TableShowView struct {
	BaseView
	table  TableModel
	window *Window
}

func NewTableShowView(table TableModel) *TableShowView {
	width := ui.TermWidth()
	height := ui.TermHeight()

	view := &TableShowView{table: table}
	dashboard := NewDashboardView(0, 0, width, 3)
	content := NewTableView(Point{0, 4}, Size{width, 8}, table.Glimpse())
	dashboard.delegate = content
	status := NewStatusBarView(Point{0, height - 3}, width, content)
	window := NewWindow(dashboard, content, status)
	view.window = window
	return view
}
