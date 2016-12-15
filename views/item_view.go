package views

import (
	"fmt"
	"strings"
)

const formatter = "fg-white,bg-blue"

type HasContent interface {
	Content(int) string
	//	EntryPoint() *Window
}

type ItemView struct {
	Object   HasContent
	width    int
	selected bool
}

func (item ItemView) Content() string {
	str := item.Object.Content(item.width)
	if item.selected {
		str = fmt.Sprintf("[%s](%s)", str, formatter)
	}
	return str
}

func (item ItemView) Match(destination string) bool {
	content := item.Object.Content(0)
	return strings.Contains(content, destination)
}

/*
func (item *ItemView) Enter() {
	window := item.Object.EntryPoint()
	if window != nil {
		nav.Push(window)
	}
}
*/
