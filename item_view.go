package main

import (
	"fmt"
	"strings"
)

const formatter = "fg-white,bg-blue"

type HasContent interface {
	Content(int) string
	EntryPoint() *Window
}

type ItemView struct {
	object   HasContent
	width    int
	selected bool
}

func (item ItemView) Content() string {
	str := item.object.Content(item.width)
	if item.selected {
		str = fmt.Sprintf("[%s](%s)", str, formatter)
	}
	return str
}

func (item ItemView) Match(destination string) bool {
	content := item.object.Content(0)
	return strings.Contains(content, destination)
}

func (item ItemView) Enter() {
	window := item.object.EntryPoint()
	if window != nil {
		nav.Push(window)
	}
}
