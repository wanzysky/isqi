package main

import termbox "github.com/nsf/termbox-go"

func print_row(left_top Point, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(left_top.x, left_top.y, c, fg, bg)
		left_top.x++
	}
}

func print_line(left_top Point, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(left_top.x, left_top.y, c, fg, bg)
		left_top.y++
	}
}

func print_content(place Point, color, background termbox.Attribute, content string, vertical bool) {
	for i, char := range content {
		if vertical {
			termbox.SetCell(place.x, place.y+i, char, color, background)
		} else {
			termbox.SetCell(place.x+i, place.y, char, color, background)
		}

	}
}
