package main

import termbox "github.com/nsf/termbox-go"
import s "strings"

type ContainerView struct {
	BaseView
	bordered bool
}

const left_devider = "┌"
const right_devider = "┐"
const left_devider_btm = "└"
const right_devider_btm = "┘"
const horizontal_border = "─"
const vertical_border = "|"
const border_color = termbox.ColorWhite
const border_bg_color = termbox.ColorBlack

func (container ContainerView) Draw() {
	location := container.location
	size := container.size

	if size.width <= 2 || size.height <= 2 {
		return
	}

	if !container.bordered {
		return
	}

	//borders
	top_frame := left_devider + s.Repeat(horizontal_border, size.width-2) + right_devider
	vertical_frame := s.Repeat(vertical_border, size.height-3)
	botom_frame := left_devider_btm + s.Repeat(horizontal_border, size.width-2) + right_devider_btm

	print_row(location, border_color, border_bg_color, top_frame)
	print_row(Point{0, size.height - 2}, border_color, border_bg_color, botom_frame)
	print_line(Point{0, 1}, border_color, border_bg_color, vertical_frame)
	print_line(Point{size.width - 1, 1}, border_color, border_bg_color, vertical_frame)
}
