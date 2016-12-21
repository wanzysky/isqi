package windows

import (
	"image"
	m "github.com/wanzysky/isqi/models"
	v "github.com/wanzysky/isqi/views"
)

type TableStructureWindow struct {
	TableShowWindow
}

func NewTableStuctureWindow(table *m.TableModel) *TableStructureWindow {
	rect := Size()
	table.SyncColumns()
	view := TableStructureWindow{}
	view.table = table
	dashboard := v.NewDashboardView(image.Rect(rect.Min.X, rect.Min.Y, rect.Max.X, 3))
	content := v.NewTableView(image.Rect(rect.Min.X, rect.Min.Y+3, rect.Max.X, rect.Max.Y-3), table.Structure())
	dashboard.Delegate = content
	status := v.NewStatusBarView(image.Rect(rect.Min.X, rect.Max.Y-3, rect.Max.X, rect.Max.Y), content)
	window := NewWindow(dashboard, content, status)
	view.table_view = content
	view.dash = dashboard
	view.status_bar = status
	view.Window = window
	return &view
}
