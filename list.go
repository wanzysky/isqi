package main

type List struct {
	items       [2]Item
	location    Point
	width       int
	line_height int
	max_height  int
}

func (list List) Draw() {
	draw_frame(list.location, Size{list.width, list.max_height + 2})
	for i, item := range list.items {
		item.location = Point{2, i*list.line_height + 2}
		item.vertical = false
		if i*list.line_height >= list.max_height {
			break
		}
		item.Draw()
	}
}
