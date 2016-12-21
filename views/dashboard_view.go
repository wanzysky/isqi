package views

import (
	ui "github.com/gizak/termui"
	termbox "github.com/nsf/termbox-go"
	"image"
	"strconv"
)

type DashState int

const (
	DashNormal = iota
	DashSearching
	DashChoosing
)

type DashboardView struct {
	BaseView
	TypingView
	Tips     string
	state    DashState
	view     *ui.Par
	Delegate Dashable
}

type Dashable interface {
	Operations() map[string]string
	SearchingTip() string
	Search(string)
	Choose(int)
	Normal()
}

func NewDashboardView(rect image.Rectangle) *DashboardView {
	dash := DashboardView{}
	dash.typing = false
	dash.rect = rect
	dash.state = DashNormal
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
	operations := dashboard.Delegate.Operations()
	tips_str := ""
	for key, op := range operations {
		tips_str += "[" + key + "] " + "[" + op + "]" + "(fg-white,bg-blue)  "
	}
	dashboard.Tips = tips_str
	return dashboard.Tips
}

func (dashboard *DashboardView) Sync() {
	if dashboard.view == nil {
		return
	}

	if dashboard.Tips == "" {
		dashboard.HelpingText()
	}

	dashboard.view.Width = dashboard.rect.Dx()
	dashboard.view.Height = dashboard.rect.Dy()
	dashboard.view.X = dashboard.rect.Min.X
	dashboard.view.Y = dashboard.rect.Min.Y

	content := dashboard.Content()
	dashboard.view.Text = content
	dashboard.cursor = dashboard.rect.Min.Add(image.Pt(len(content)+1, 1))
}

func (dashboard *DashboardView) Content() string {
	var str string
	switch dashboard.state {
	case DashNormal:
		str = dashboard.Tips
	case DashChoosing:
		str = "[Choosing by number] " + dashboard.value
	case DashSearching:
		str = dashboard.Delegate.SearchingTip() + dashboard.value
	}
	return str
}

func (dashboard *DashboardView) Draw() {
	if dashboard.typing {
		termbox.SetCursor(dashboard.cursor.X, dashboard.cursor.Y)
	} else {
		termbox.HideCursor()
	}

	view := dashboard.View()
	ui.Render(view)
}

func (dashboard *DashboardView) Display() {
	dashboard.Draw()
}

func (dashboard *DashboardView) Clear() {
	ui.DefaultEvtStream.ResetHandlers()
	ui.Clear()
}

func (dashboard *DashboardView) Key(key string) bool {
	switch key {
	case "s":
		if dashboard.typing {
			dashboard.Type("s")
		} else {
			dashboard.Searching()
		}
		return true
	case "c":
		if dashboard.typing {
			dashboard.Type("c")
		} else {
			dashboard.Choosing()
		}
		return true
	case "C-8":
		if dashboard.typing {
			dashboard.Delete()
		}
		return true
	default:
		if dashboard.typing {
			dashboard.Type(key)
			return true
		}
	}
	return false
}

func (dashboard *DashboardView) Searching() {
	//dashboard.Tips = dashboard.Delegate.SearchingTip()
	dashboard.typing = true
	dashboard.state = DashSearching
	dashboard.value = ""
	dashboard.Sync()
	dashboard.Draw()
}

func (dashboard *DashboardView) Choosing() {
	dashboard.typing = true
	dashboard.state = DashChoosing
	dashboard.value = ""
	dashboard.Sync()
	dashboard.Draw()
}

func (dashboard *DashboardView) Choose() {
	index, err := strconv.Atoi(dashboard.value)
	if err != nil {
		return
	}
	dashboard.typing = true
	dashboard.state = DashChoosing
	dashboard.Delegate.Choose(index)
}

func (dashboard *DashboardView) Type(word string) {
	if len(word) > 1 {
		return
	}
	dashboard.value += word
	dashboard.Sync()
	dashboard.Draw()
	if dashboard.state == DashSearching {
		dashboard.Delegate.Search(dashboard.value)
	} else if dashboard.state == DashChoosing {
		dashboard.Choose()
	}
}

func (dashboard *DashboardView) Delete() {
	if len(dashboard.value) < 1 {
		return
	}

	dashboard.value = dashboard.value[:len(dashboard.value)-1]
	dashboard.Sync()
	dashboard.Draw()
	if dashboard.state == DashSearching {
		dashboard.Delegate.Search(dashboard.value)
	} else if dashboard.state == DashChoosing {
		dashboard.Choose()
	}
}

func (dashboard *DashboardView) Escape() bool {
	if dashboard.typing {
		dashboard.Normal()
		dashboard.Delegate.Normal()
		return true
	} else {
		return false
	}
}

func (dashboard *DashboardView) Normal() {
	dashboard.typing = false
	dashboard.state = DashNormal
	dashboard.Sync()
	dashboard.Draw()
}
