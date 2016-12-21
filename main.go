package main

import (
	ui "github.com/gizak/termui"
	adpt "isqi/adapters"
	wd "isqi/windows"
)

func main() {
	config := Config()
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	window := config.Connect()
	defer adpt.Conn.Close()

	wd.NewNavigatoin(window)
	ui.Loop()
}
