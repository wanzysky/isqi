package main

import ui "github.com/gizak/termui"

type TableView struct {
	BaseView
	datasource  [][]string
	word_length []int
	topleft     Point
	view        *ui.List
}

func NewTableView(position Point, size Size, datasource [][]string) *TableView {
	table_view := TableView{datasource: datasource}
	table_view.location = position
	table_view.size = size
	if len(datasource) > 0 {
		table_view.word_length = make([]int, len(datasource[0]))
		for _, row := range datasource {
			for i, cell := range row {
				length := len(cell)
				if length > table_view.word_length[i] {
					table_view.word_length[i] = length
				}
			}
		}
	}
	return &table_view
}

func (tableview *TableView) Serialize() {

}

func (tableview *TableView) Sync() {

}

func (tableview *TableView) Draw() {
}
