package windows

import (
	ui "github.com/gizak/termui"
	v "github.com/wanzysky/isqi/views"
)

type ColumnDetailWindow struct {
	*Window
	col *v.CollectionView
}

func NewColumnDetailWindow(headers, contents []string) *ColumnDetailWindow {
	rect := Size()
	window := ColumnDetailWindow{}
	collection_view := v.NewCollectionView(rect, headers, contents)
	base := NewWindow(collection_view)
	window.Window = base
	window.col = collection_view
	return &window
}

func (window *ColumnDetailWindow) Display() {
	view := window.views
	for view != nil {
		view.Display()
		view = view.next
	}

	window.Listening()
}

func (window *ColumnDetailWindow) Listening() {
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/<escape>", func(ui.Event) {
		Nav.Back()
	})

	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		window.col.Down()
	})

	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		window.col.Up()
	})

}
