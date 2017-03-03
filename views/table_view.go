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
	page            int
	per_page        int
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
	table_view.view.Separator = false
	table_view.current = 0
	table_view.offset = 0
	table_view.per_page = rect.Dy() - 2
	table_view.Serialize()

	return &table_view
}

func (tableview *TableView) Serialize() {
	if len(tableview.datasource) <= 0 {
		tableview.rows = [][]string{}
		return
	}
	datasource := tableview.CurrentPage()
	tableview.max_lengh = make([]int, len(datasource[0]))
	tableview.rows = make([][]string, len(datasource))

	if len(datasource) <= 0 {
		return
	}

	for j, row := range tableview.CurrentPage() {
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

func (tableview *TableView) Update(source [][]string) {
	tableview.datasource = source
	tableview.view = &ui.Table{}
	tableview.view.Width = tableview.rect.Dx()
	tableview.view.Height = tableview.rect.Dy()
	tableview.view.X = tableview.rect.Min.X
	tableview.view.Y = tableview.rect.Min.Y
	tableview.view.FgColor = ui.ColorWhite
	tableview.view.BgColor = ui.ColorDefault
	tableview.view.TextAlign = ui.AlignCenter
	tableview.view.Separator = false
	tableview.current = 0
	tableview.offset = 0
	tableview.per_page = tableview.rect.Dy() - 2
	tableview.Serialize()
	tableview.Display()
}

func (tableview *TableView) ReDraw() {
	ui.ClearArea(tableview.rect, DEFAULT_BG_COLOR)
	ui.Render(tableview.view)
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

func (tableview *TableView) CurrentPage() [][]string {
	start := tableview.per_page * tableview.page
	ending := start + tableview.per_page
	length := len(tableview.datasource)
	if ending > length {
		ending = length
	}
	return tableview.datasource[start:ending]
}

func (tableview *TableView) PageTo(page int) {
	if page < 0 {
		return
	}
	length := len(tableview.datasource)
	if page*tableview.per_page > length {
		return
	}

	tableview.page = page
	tableview.rows = tableview.CurrentPage()
}

func (tableview *TableView) PageUp() {
	tableview.PageTo(tableview.page - 1)
	tableview.Sync()
	tableview.ReDraw()
}

func (tableview *TableView) PageDown() {
	tableview.PageTo(tableview.page + 1)
	tableview.Sync()
	tableview.ReDraw()
}

func (tableview *TableView) Left() {
	if tableview.offset == 0 {
		return
	}
	tableview.offset -= 1
	tableview.Sync()
	tableview.ReDraw()
}

func (tableview *TableView) Right() {
	if len(tableview.rows[0])-tableview.offset == 1 {
		return
	}
	tableview.offset += 1
	tableview.Sync()
	tableview.ReDraw()
}

func (tableview *TableView) Current() ([]string, []string) {
	if len(tableview.datasource) == 0 {
		return []string{}, []string{}
	}

	return tableview.datasource[0], tableview.rows[tableview.current]
}

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
		"C-f/b": "Page up/down",
		"C-c":   "Quit",
		"Enter": "Console",
		"d":     "Detail",
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
