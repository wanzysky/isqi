package windows

import (
	ui "github.com/gizak/termui"
	"image"
)

type Drawable interface {
	Display()
	Clear()
}

type View struct {
	Drawable
	next *View
}

type Window struct {
	views      *View
	rect       image.Rectangle
	evt_stream *ui.EvtStream
	active     bool
}

func NewWindow(views ...Drawable) *Window {
	var window Window
	window.rect = ui.TermRect()
	//window.evt_stream = ui.NewEvtStream()
	//window.evt_stream.Init()
	if len(views) > 0 {
		var current *View
		for _, v := range views {
			view := View{Drawable: v}
			if window.views == nil {
				window.views = &view
				current = window.views
			} else {
				current.next = &view
				current = current.next
			}
		}
	}
	return &window
}

func Size() image.Rectangle {
	return ui.TermRect()
}

func (window *Window) Clear() {
	ui.Clear()
	ui.DefaultEvtStream.ResetHandlers()
	window.active = false
}
