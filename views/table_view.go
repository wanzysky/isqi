package views

import (
	"fmt"

	ui "github.com/gizak/termui"
	"image"
)

type TableView struct {
	BaseView
	datasource      [][]string
	rows            [][]string
	max_lengh       []int
	topleft         image.Point
	current         int
	offset          int
	view            *ui.Table
	status_delegate *StatusBarView
}

const MAX_CELL_WIDTH = 30
const ACTIVE_BG_COLOR = ui.ColorBlue
const DEFAULT_BG_COLOR = ui.ColorDefault

func NewTableView(rect image.Rectangle, datasource [][]string) *TableView {
	table_view := TableView{datasource: datasource}
	table_view.rect = rect
	table_view.view = &ui.Table{}
	table_view.view.Width = table_view.rect.Dx()
	table_view.view.Height = table_view.rect.Dy()
	table_view.view.X = table_view.rect.Min.X
	table_view.view.Y = table_view.rect.Min.Y
	table_view.view.FgColor = ui.ColorWhite
	table_view.view.BgColor = ui.ColorDefault
	table_view.view.TextAlign = ui.AlignCenter
	table_view.view.Seperator = false
	table_view.current = 0
	table_view.offset = 0
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

func (tableview *TableView) OffsetRows() (rows [][]string) {
	for _, row := range tableview.rows {
		rows = append(rows, row[tableview.offset:])
	}
	return
}

func (tableview *TableView) Sync() {
	tableview.view.Rows = tableview.OffsetRows()
	tableview.view.Analysis()
	tableview.view.SetSize()
	tableview.view.BgColors[tableview.current] = ACTIVE_BG_COLOR
}

func (tableview *TableView) Draw() {
	ui.Render(tableview.view)
}

// Events
func (tableview *TableView) Up() {
	current := tableview.current - 1
	if current < 0 {
		current = len(tableview.rows) - 1
	}
	tableview.Select(current)
	tableview.Sync()
	tableview.Draw()
}

func (tableview *TableView) Down() {
	current := (tableview.current + 1) % (len(tableview.rows))
	tableview.Select(current)
	tableview.Sync()
	tableview.Draw()
}

func (tableview *TableView) Select(new_current int) {
	tableview.view.BgColors[tableview.current] = DEFAULT_BG_COLOR
	tableview.current = new_current
}

func (tableview *TableView) Left() {
	if tableview.offset == 0 {
		return
	}
	tableview.offset -= 1
	ui.ClearArea(image.Rect(tableview.view.X, tableview.view.Y, tableview.view.X+tableview.view.Width, tableview.view.Y+tableview.view.Height), DEFAULT_BG_COLOR)
	tableview.Sync()
	tableview.Draw()
}

func (tableview *TableView) Right() {
	if len(tableview.rows[0])-tableview.offset == 1 {
		return
	}
	tableview.offset += 1
	ui.ClearArea(image.Rect(tableview.view.X, tableview.view.Y, tableview.view.X+tableview.view.Width, tableview.view.Y+tableview.view.Height), DEFAULT_BG_COLOR)
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

func (tableview TableView) Normal() {
}

func (tableview TableView) Choose(int) {
}

// Drawable Interface
func (tableview *TableView) Display() {
	tableview.Sync()
	tableview.Draw()
}

func (tableview TableView) Clear() {
	ui.DefaultEvtStream.ResetHandlers()
	ui.Clear()
}

func (tableview *TableView) Notice(content string) {
	if tableview.status_delegate != nil {
		tableview.status_delegate.Notice(content)
	}
}
