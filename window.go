package main

type Drawable interface {
	Display()
	Clear()
}

type Window struct {
	views  []Drawable
	active bool
}

func NewWindow(views ...Drawable) *Window {
	var window Window
	window.views = views
	return &window
}

func (window *Window) Clear() {
	for _, view := range window.views {
		view.Clear()
	}
	window.active = false
}
func (window *Window) Display() {
	for _, view := range window.views {
		view.Display()
	}
	window.active = true
}
