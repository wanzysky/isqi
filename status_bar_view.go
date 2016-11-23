package main

import ui "github.com/wanzysky/termui"

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

func NewStatusBarView(position Point, width int, delegate Statusable) *StatusBarView {
	status_bar := &StatusBarView{}
	status_bar.location = position
	status_bar.size = Size{width, 3}
	status_bar.view = ui.NewGauge()
	status_bar.view.X = position.x
	status_bar.view.Y = position.y
	status_bar.view.Width = width
	status_bar.view.Height = 3
	return status_bar
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
