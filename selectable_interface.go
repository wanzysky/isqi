package main

type Selectable interface {
	Current()
	Up()
	Down()
	Selected()
	Typing()
}
