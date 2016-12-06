package main

import (
	"fmt"

	ui "github.com/gizak/termui"
)

type TableView struct {
	BaseView
	datasource [][]string
	rows       [][]string
	max_lengh  []int
	topleft    Point
	current    int
	view       *ui.Table
}

const MAX_CELL_WIDTH = 30
const ACTIVE_BG_COLOR = ui.ColorBlue

func NewTableView(position Point, size Size, datasource [][]string) *TableView {
	table_view := TableView{datasource: datasource}
	table_view.location = position
	table_view.size = size
	table_view.view = &ui.Table{}
	table_view.view.Width = table_view.size.width
	table_view.view.Height = table_view.size.height
	table_view.view.X = table_view.location.x
	table_view.view.Y = table_view.location.y
	table_view.view.FgColor = ui.ColorWhite
	table_view.view.BgColor = ui.ColorDefault
	table_view.view.TextAlign = ui.AlignCenter
	table_view.view.Seperator = false
	table_view.current = 0
	table_view.Serialize()

	return &table_view
}

func (tableview *TableView) Serialize() {
	if len(tableview.datasource) <= 0 {
		tableview.rows = [][]string{}
		return
	}
	tableview.max_lengh = make([]int, len(tableview.datasource[0]))
	tableview.rows = make([][]string, len(tableview.datasource))

	if len(tableview.datasource) <= 0 {
		return
	}

	for j, row := range tableview.datasource {
		tableview.rows[j] = make([]string, len(row))
		for i, cell := range row {
			length := len(cell)
			max_cell_width := length
			if max_cell_width > MAX_CELL_WIDTH {
				max_cell_width = MAX_CELL_WIDTH
			}

			if max_cell_width > tableview.max_lengh[i] {
				tableview.max_lengh[i] = max_cell_width
			}
			tableview.rows[j][i] = ""
			if length >= max_cell_width {
				tableview.rows[j][i] = cell[:max_cell_width]
			} else {
				tableview.rows[j][i] = cell[:length]
			}
		}
	}
}

func (tableview *TableView) Sync() {
	tableview.view.Rows = tableview.rows
	tableview.view.Analysis()
	tableview.view.SetSize()
	tableview.view.BgColors[tableview.current] = ACTIVE_BG_COLOR
}

func (tableview *TableView) Listening() {
	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		tableview.Down()
	})

	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		tableview.Up()
	})

	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		tableview.Enter()
	})
}

func (tableview *TableView) Draw() {
	ui.Render(tableview.view)
}

// Events
func (tableview *TableView) Up() {
	tableview.current -= 0
	if tableview.current < 0 {
		tableview.current = len(tableview.rows)
	}
	tableview.Sync()
	tableview.Draw()
}

func (tableview *TableView) Down() {
	tableview.current = (tableview.current + 1) % (len(tableview.rows))
	tableview.Sync()
	tableview.Draw()
}

func (tableview *TableView) Enter() {}

// Statusable Interface
func (tableview TableView) Loading(percent int) string {
	return fmt.Sprintf("Executing:%d%%", percent)
}

func (tableview TableView) Succeed() string {
	return "Execution finished"
}

func (tableview TableView) Failed() string {
	return "There are errors! "
}

// Dashable Interface
func (tableview TableView) Operations() map[string]string {
	operatios := map[string]string{
		"s":     "Search",
		"c":     "Quick Choose",
		"d":     "Database Detail",
		"C-c":   "Quit",
		"Enter": "Edit",
	}
	return operatios
}

func (tableview TableView) SearchingTip() string {
	return "[Searching]"
}

func (tableview TableView) Search(keyword string) {
}

func (tableview TableView) Choose(int) {
}

// Drawable Interface
func (tableview TableView) Display() {
	tableview.Sync()
	tableview.Draw()
}

func (tableview TableView) Clear() {
	ui.DefaultEvtStream.ResetHandlers()
	ui.Clear()
}
