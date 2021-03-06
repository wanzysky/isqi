package views

import (
	ui "github.com/gizak/termui"
	"image"
)

type StatusBarState uint8

const (
	normal = iota
	loading
	succeed
	failed
)

type StatusBarView struct {
	BaseView
	percent  int
	label    string
	view     *ui.Gauge
	state    StatusBarState
	delegate Statusable
}

type Statusable interface {
	Loading(int) string
	Succeed() string
	Failed() string
}

func NewStatusBarView(rect image.Rectangle, delegate Statusable) *StatusBarView {
	status_bar := &StatusBarView{}
	status_bar.rect = rect
	status_bar.view = ui.NewGauge()
	status_bar.view.X = rect.Min.X
	status_bar.view.Y = rect.Min.Y
	status_bar.view.Width = rect.Dx()
	status_bar.view.Height = 3
	return status_bar
}

func (bar *StatusBarView) Success(message string) {
	bar.state = succeed
	bar.label = message
}

func (bar *StatusBarView) Notice(message string) {
	bar.label = message
	bar.Sync()
	bar.Draw()
}

func (status_bar *StatusBarView) Sync() {
	status_bar.view.Percent = status_bar.percent
	status_bar.view.Label = status_bar.label
	status_bar.view.PercentColor = ui.ColorWhite
	status_bar.view.PercentColorHighlighted = ui.ColorBlack

	switch status_bar.state {
	case normal:
		status_bar.view.Percent = 0
	case loading:
		status_bar.view.BarColor = ui.ColorWhite
	case succeed:
		status_bar.view.Percent = 100
		status_bar.view.BarColor = ui.ColorGreen
	case failed:
		status_bar.view.Percent = 100
		status_bar.view.BarColor = ui.ColorYellow
	}
}

func (status_bar *StatusBarView) Draw() {
	ui.Render(status_bar.view)
}

// Drawable Interface
func (status_bar StatusBarView) Display() {
	status_bar.Sync()
	status_bar.Draw()
}

func (status_bar StatusBarView) Clear() {
	ui.DefaultEvtStream.ResetHandlers()
	ui.Clear()
}
