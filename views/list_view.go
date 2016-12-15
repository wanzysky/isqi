package views

import (
	"fmt"

	ui "github.com/gizak/termui"
	"image"
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

func NewListView(rect image.Rectangle, title string, subs []ItemView) *ListView {
	list := ListView{title: title, datasource: subs}
	list.rect = rect
	list.title = title
	list.view = list.View()
	for _, item := range subs {
		item.width = rect.Dx() - 2
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
		view.X = list.rect.Min.X
		view.Y = list.rect.Min.Y
		view.Width = list.rect.Dx()
		view.Height = list.rect.Dy()
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

func (list *ListView) Current() *ItemView {
	return list.current
}

// Dashable Interface
func (list *ListView) Search(content string) {
	var new_items []ItemView
	for _, data := range list.datasource {
		if data.Match(content) {
			new_items = append(new_items, data)
		}
	}
	list.items = new_items
	list.Select(0)
	list.Sync()
	list.Draw()
}

func (list *ListView) Display() {
	list.Sync()
	list.Draw()
}

func (list *ListView) Clear() {
}

func (list ListView) Operations() map[string]string {
	operatios := map[string]string{
		"s":     "Search",
		"c":     "Quick Choose",
		"d":     "Database Detail",
		"C-c":   "Quit",
		"Enter": "Use",
	}
	return operatios
}

func (list *ListView) Choose(index int) {
	list.Select(index)
	list.Sync()
	list.Draw()
}

func (list *ListView) Normal() {
	list.items = list.datasource
	list.Sync()
	list.Draw()
}

func (list ListView) SearchingTip() string {
	return "[Searching In List] "
}
