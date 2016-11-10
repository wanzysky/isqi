package main

import (
	"fmt"

	ui "github.com/gizak/termui"
)

const content_color = ui.ColorYellow
const content_bg_color = ui.ColorBlack
const active_bg_color = ui.ColorBlue

type ListView struct {
	BaseView
	datasource    []ItemView
	items         []ItemView
	title         string
	current_index int
	current       *ItemView
	view          *ui.List
}

func NewListView(x, y int, width, height int, title string, subs []ItemView) *ListView {
	list := ListView{title: title, datasource: subs}
	list.location = Point{x, y}
	list.size = Size{width, height}
	list.title = title
	list.view = list.View()
	for _, item := range subs {
		item.width = width - 2
	}

	list.items = subs
	if len(subs) > 0 {
		list.current_index = 0
		list.current = &list.items[0]
		list.current.selected = true
	}
	return &list
}

func (list *ListView) Select(index int) {
	length := len(list.items)
	if !(length > 0) {
		return
	}

	if index < 0 {
		index = length - 1
	} else if index >= length {
		index %= length
	}

	list.current.selected = false
	list.current_index = index
	list.current = &list.items[list.current_index]
	list.current.selected = true
}

func (list *ListView) Up() {
	list.Select(list.current_index - 1)
}

func (list *ListView) Down() {
	list.Select(list.current_index + 1)
}

func (list ListView) View() *ui.List {
	if list.view == nil {
		view := ui.NewList()
		list.view = view
		view.ItemFgColor = content_color
		view.X = list.location.x
		view.Y = list.location.y
		view.Width = list.size.width
		view.Height = list.size.height
		view.BorderLabel = list.title
	}
	return list.view
}

func (list ListView) Sync() {
	var strs []string
	for i, item := range list.items {
		str := fmt.Sprintf("[%d] %s", i, item.Content())
		strs = append(strs, str)
	}
	list.view.Items = strs
}

func (list ListView) Draw() {
	view := list.View()
	ui.Render(view)
}

func (list ListView) Listening() {
	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		list.Down()
		list.Sync()
		ui.Render(list.view)
	})

	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		list.Up()
		list.Sync()
		ui.Render(list.view)
	})

	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		list.current.Enter()
	})
}

func (list ListView) Search(content string) {
	var new_items []ItemView
	for _, data := range list.datasource {
		if data.Match(content) {
			new_items = append(new_items, data)
		}
	}
	list.items = new_items
	list.Sync()
	list.Draw()
}

func (list ListView) Display() {
	list.Sync()
	list.Draw()
	list.Listening()
}

func (list ListView) Clear() {
	ui.DefaultEvtStream.ResetHandlers()
	ui.Clear()
}

func (list ListView) Choose(index int) {

}

func (list ListView) SearchingTip() string {
	return "[Searching In List] "
}
