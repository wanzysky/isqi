package main

import (
	ui "github.com/gizak/termui"
	"go-webterm"
	adpt "isqi/adapters"
	wd "isqi/windows"
)

func main() {
	go func() { debuger.ListenAndServe() }()

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	config := Config()
	window := config.Connect()
	defer adpt.Conn.Close()

	wd.NewNavigatoin(window)
	ui.Loop()
}
