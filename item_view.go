package main

import termbox "github.com/nsf/termbox-go"

var content_color = termbox.ColorWhite
var content_bg_color = termbox.ColorBlack

type ItemView struct {
	BaseView
	content  string
	selected bool
}

func (item Item) Draw() {
	print_content(item.location, content_color, content_bg_color, item.content, item.vertical)
}
