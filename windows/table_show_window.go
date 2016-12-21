package windows

import (
	"fmt"
	ui "github.com/gizak/termui"
	"image"
	m "github.com/wanzysky/isqi/models"
	v "github.com/wanzysky/isqi/views"
)

type TableShowWindow struct {
	table *m.TableModel
	*Window
	table_view *v.TableView
	dash       *v.DashboardView
	status_bar *v.StatusBarView
}

func NewTableShowWindow(table *m.TableModel) *TableShowWindow {
	rect := Size()
	view := TableShowWindow{table: table}
	dashboard := v.NewDashboardView(image.Rect(rect.Min.X, rect.Min.Y, rect.Max.X, 3))
	content := v.NewTableView(image.Rect(rect.Min.X, rect.Min.Y+3, rect.Max.X, rect.Max.Y-3), table.Glimpse())
	dashboard.Delegate = content
	status := v.NewStatusBarView(image.Rect(rect.Min.X, rect.Max.Y-3, rect.Max.X, rect.Max.Y), content)
	count := table.Statistic()
	shown := 100
	if count < shown {
		shown = count
	}
	status.Success(fmt.Sprintf("Rows 1 - %d of %d from table", shown, count))
	window := NewWindow(dashboard, content, status)
	view.table_view = content
	view.dash = dashboard
	view.status_bar = status
	view.Window = window
	return &view
}

func (window *TableShowWindow) Listening() {
	tableview := window.table_view
	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		tableview.Down()
	})

	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		tableview.Up()
	})

	ui.Handle("/sys/kbd/<left>", func(ui.Event) {
		tableview.Left()
	})

	ui.Handle("/sys/kbd/<right>", func(ui.Event) {
		tableview.Right()
	})

	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		window.Enter()
	})

	ui.Handle("/sys/kbd/d", func(ui.Event) {
		window.Detail(tableview.Current())
	})

	ui.Handle("/sys/kbd/C-f", func(ui.Event) {
		tableview.PageDown()
	})

	ui.Handle("/sys/kbd/C-b", func(ui.Event) {
		tableview.PageUp()
	})

	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/<escape>", func(ui.Event) {
		if !window.dash.Escape() {
			Nav.Back()
		}
	})

	ui.Handle("/sys/kbd/", func(e ui.Event) {
		window.dash.Key(e.Data.(ui.EvtKbd).KeyStr)
	})

}

func (window *TableShowWindow) Display() {
	view := window.views
	for view != nil {
		view.Display()
		view = view.next
	}
	window.active = true
	window.Listening()
}

func (window *TableShowWindow) Enter() {
	Nav.Push(NewConsoleWindow())
}

func (window *TableShowWindow) Detail(headers []string, contents []string) {
	Nav.Push(NewColumnDetailWindow(headers, contents))
}
