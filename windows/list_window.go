package windows

import (
	ui "github.com/gizak/termui"
	m "isqi/models"
	v "isqi/views"
)

type ListWindow struct {
	*Window
	list *v.ListView
	dash *v.DashboardView
}

func NewListWindow(list *v.ListView, dash *v.DashboardView) *ListWindow {
	window := ListWindow{Window: NewWindow(list, dash)}
	window.list = list
	window.dash = dash
	return &window
}

func (window *ListWindow) Listening() {
	list := window.list
	dashboard := window.dash

	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		list.Down()
		list.Display()
	})

	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		list.Up()
		list.Display()
	})

	ui.Handle("/sys/kbd/C-f", func(ui.Event) {
		list.PageDown()
		list.Display()
	})

	ui.Handle("/sys/kbd/C-b", func(ui.Event) {
		list.PageUp()
		list.Display()
	})

	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		window.Enter(list.Current())
	})

	ui.Handle("/sys/kbd/d", func(ui.Event) {
		if !window.dash.Key("d") {
			window.Detail(list.Current())
		}
	})

	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/<escape>", func(ui.Event) {
		if !dashboard.Escape() {
			Nav.Back()
		}
	})

	ui.Handle("/sys/kbd/", func(e ui.Event) {
		dashboard.Key(e.Data.(ui.EvtKbd).KeyStr)
	})
}
func (window *ListWindow) Detail(item *v.ItemView) {
	if table, ok := item.Object.(*m.TableModel); ok {
		Nav.Push(NewTableStuctureWindow(table))
	}
}

func (window *ListWindow) Enter(item *v.ItemView) {
	if db, ok := item.Object.(m.DatabaseModel); ok {
		Nav.Push(NewTableIndexWindow(&db))
	}

	if table, ok := item.Object.(*m.TableModel); ok {
		Nav.Push(NewTableShowWindow(table))
	}
}

func (window *ListWindow) Clear() {
	ui.Clear()
	ui.DefaultEvtStream.ResetHandlers()
}

func (window *ListWindow) Display() {
	view := window.views
	for view != nil {
		view.Display()
		view = view.next
	}
	window.active = true
	window.Listening()
}
