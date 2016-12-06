package main

import (
	termbox "github.com/nsf/termbox-go"
	ui "github.com/gizak/termui"
)

type DashboardView struct {
	BaseView
	tips     string
	typing   bool
	value    string
	cursor   Point
	view     *ui.Par
	delegate Dashable
}

type Dashable interface {
	Operations() map[string]string
	SearchingTip() string
	Search(string)
	Choose(int)
}

func NewDashboardView(x, y, width, height int) *DashboardView {
	dash := DashboardView{typing: false}
	dash.location = Point{x, y}
	dash.size = Size{width, height}
	dash.value = ""
	return &dash
}

func (dashboard *DashboardView) View() *ui.Par {
	if dashboard.view == nil {
		view := ui.NewPar(dashboard.Content())
		dashboard.view = view
		dashboard.Sync()
	}
	return dashboard.view
}

func (dashboard *DashboardView) HelpingText() string {
	operations := dashboard.delegate.Operations()
	tips_str := ""
	for key, op := range operations {
		tips_str += "[" + key + "] " + "[" + op + "]" + "(fg-white,bg-blue)  "
	}
	dashboard.tips = tips_str
	return dashboard.tips
}

func (dashboard *DashboardView) Sync() {
	if dashboard.view == nil {
		return
	}

	if dashboard.tips == "" {
		dashboard.HelpingText()
	}

	dashboard.view.Width = dashboard.size.width
	dashboard.view.Height = dashboard.size.height
	dashboard.view.X = dashboard.location.x
	dashboard.view.Y = dashboard.location.y

	content := dashboard.Content()
	dashboard.view.Text = content
	dashboard.cursor = Point{x: dashboard.location.x + len(content) + 1, y: dashboard.location.y + 1}
}

func (dashboard *DashboardView) Content() string {
	var str string
	if dashboard.typing {
		str = dashboard.delegate.SearchingTip() + dashboard.value
	} else {
		str = dashboard.tips
	}
	return str
}

func (dashboard *DashboardView) Draw() {
	if dashboard.typing {
		termbox.SetCursor(dashboard.cursor.x, dashboard.cursor.y)
	} else {
		termbox.HideCursor()
	}

	view := dashboard.View()
	ui.Render(view)
}

func (dashboard DashboardView) Display() {
	dashboard.Draw()
	dashboard.Listening()
}

func (dashboard DashboardView) Clear() {
	ui.DefaultEvtStream.ResetHandlers()
	ui.Clear()
}

func (dashboard *DashboardView) Searching() {
	//dashboard.tips = dashboard.delegate.SearchingTip()
	dashboard.typing = true
	dashboard.value = ""
	dashboard.Sync()
	dashboard.Draw()
}

func (dashboard *DashboardView) Type(word string) {
	//if len(word) > 1 {
	//	return
	//}
	dashboard.value += word
	dashboard.Sync()
	dashboard.Draw()
	dashboard.delegate.Search(dashboard.value)
}

func (dashboard *DashboardView) Delete() {
	if len(dashboard.value) < 1 {
		return
	}

	dashboard.value = dashboard.value[:len(dashboard.value)-1]
	dashboard.Sync()
	dashboard.Draw()
	dashboard.delegate.Search(dashboard.value)
}

func (dashboard *DashboardView) Listening() {
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/<escape>", func(ui.Event) {
		dashboard.Escape()
	})

	ui.Handle("/sys/kbd/s", func(ui.Event) {
		if dashboard.typing {
			dashboard.Type("s")
		} else {
			dashboard.Searching()
		}
	})

	ui.Handle("/sys/kbd/C-8", func(ui.Event) {
		if dashboard.typing {
			dashboard.Delete()
		}
	})

	ui.Handle("/sys/kbd/", func(e ui.Event) {
		if dashboard.typing {
			dashboard.Type(e.Data.(ui.EvtKbd).KeyStr)
		} else {
			return
		}
	})
}

func (dashboard *DashboardView) Escape() {
	if dashboard.typing {
		dashboard.Normal()
	} else {
		nav.Back()
	}
}

func (dashboard *DashboardView) Normal() {
	dashboard.typing = false
	dashboard.Sync()
	dashboard.Draw()
}
