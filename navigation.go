package main

import ui "github.com/gizak/termui"

type Pair struct {
	index int
	value *Window
}

type Navigation struct {
	windows []*Window
	current Pair
}

func NewNavigatoin(window *Window) *Navigation {
	var nav Navigation
	nav.windows = make([]*Window, 0)
	nav.windows = append(nav.windows, window)
	nav.current = Pair{0, window}
	window.Display()
	return &nav
}

func (nav *Navigation) Push(window *Window) *Navigation {
	current_index := nav.current.index
	nav.current.value.Clear()
	nav.windows = append(nav.windows[:current_index+1], window)
	nav.current.index += 1
	nav.current.value = window
	nav.current.value.Display()
	return nav
}

func (nav *Navigation) Back() *Navigation {
	if nav.current.index > 0 {
		nav.current.index -= 1
		nav.current.value.Clear()
		nav.current.value = nav.windows[nav.current.index]
		nav.current.value.Display()
	}
	return nav
}

func (nav *Navigation) Position() Point {
	return Point{0, 0}
}

func (nav *Navigation) Size() Size {
	return Size{ui.TermWidth(), ui.TermHeight()}
}
