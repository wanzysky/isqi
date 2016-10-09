package main

import "github.com/nsf/termbox-go"

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputAlt)
	width, height := termbox.Size()
	draw_frame(Point{0, 0}, Size{width, height})
	databases := [2]Item{Item{content: "db 1"}, Item{content: "db 2"}}
	list := List{items: databases, location: Point{1, 1}, width: width / 4, line_height: 1, max_height: height - 2}
	list.Draw()
	termbox.Flush()
loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlC {
				break loop
			}
		}
	}
}
