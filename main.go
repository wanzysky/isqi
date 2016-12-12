package main

import (
	"database/sql"

	ui "github.com/gizak/termui"
	"go-webterm"
)

var nav *Navigation
var connection *sql.DB
var adapter Adapter

func main() {
	go func() { debuger.ListenAndServe() }()

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	config := Config()
	window := config.Connect()
	defer connection.Close()

	nav = NewNavigatoin(window)

	ui.Loop()
}
