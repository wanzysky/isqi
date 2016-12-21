package windows

import (
	"fmt"
	ui "github.com/gizak/termui"
	"image"
	a "github.com/wanzysky/isqi/adapters"
	v "github.com/wanzysky/isqi/views"
)

type ConsoleWindow struct {
	*Window
	console_view *v.ConsoleView
	table_view   *v.TableView
	status_bar   *v.StatusBarView
	typing       bool
}

func NewConsoleWindow() *ConsoleWindow {
	rect := Size()
	window := ConsoleWindow{}

	console := v.NewConsoleView(image.Rect(rect.Min.X, rect.Min.Y, rect.Max.X, 5))
	window.console_view = console

	tableview := v.NewTableView(image.Rect(rect.Min.X, rect.Min.Y+5, rect.Max.X, rect.Max.Y-3), [][]string{[]string{"C-r to run sql"}})
	status := v.NewStatusBarView(image.Rect(rect.Min.X, rect.Max.Y-3, rect.Max.X, rect.Max.Y), tableview)
	window.table_view = tableview
	window.status_bar = status
	window.typing = true

	base := NewWindow(console, tableview, status)
	window.Window = base
	window.rect = rect
	return &window
}

func (window *ConsoleWindow) Display() {
	view := window.views
	for view != nil {
		view.Display()
		view = view.next
	}

	window.active = true
	window.Listening()
}

func (window *ConsoleWindow) Listening() {
	ui.Handle("/sys/kbd/C-r", func(ui.Event) {
		window.Exec(window.console_view.Val())
	})

	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/<escape>", func(ui.Event) {
		Nav.Back()
	})

	ui.Handle("/sys/kbd/", func(e ui.Event) {
		window.console_view.Key(e.Data.(ui.EvtKbd).KeyStr)
	})

	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		if window.typing {
			window.console_view.Key("<enter>")
		} else {
			window.console_view.Continue()
			window.typing = true
		}
	})

	tableview := window.table_view
	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		tableview.Down()
	})

	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		tableview.Up()
	})

	ui.Handle("/sys/kbd/<left>", func(ui.Event) {
		if window.typing {
			window.console_view.Key("<left>")
		} else {
			tableview.Left()
		}
	})

	ui.Handle("/sys/kbd/<right>", func(ui.Event) {
		if window.typing {
			window.console_view.Key("<right>")
		} else {
			tableview.Right()
		}
	})

	ui.Handle("/sys/kbd/C-f", func(ui.Event) {
		if !window.typing {
			tableview.PageDown()
		}
	})

	ui.Handle("/sys/kbd/C-b", func(ui.Event) {
		if !window.typing {
			tableview.PageUp()
		}
	})

}

func (window *ConsoleWindow) Exec(query string) {
	table, err := a.Adpt.Execute(query)
	if err == nil {
		window.console_view.Stop()
		window.typing = false
		window.table_view.Update(table)
		window.status_bar.Notice(fmt.Sprintf("No errors; %d rows listed, Press Enter to query", len(table)))
	} else {
		window.status_bar.Notice(err.Error())
	}
}
