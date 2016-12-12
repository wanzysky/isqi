package main

import (
	"fmt"
	ui "github.com/gizak/termui"
	"os"
)

type DbConf struct {
	host     string
	username string
	passwd   string
	database string
	port     string
}

type ConfMode int

const (
	ConfModeNormal = iota
	ConfModeUsage
	ConfModeRails
	ConfModeConfigFile
	ConfModeEnterPassword
	ConfModeInvalid
)

type Configuration struct {
	DbConf
	mode           ConfMode
	sourceFilePath string
}

type MachineState int

const (
	MachineStateContent = iota
	MachineStateHeading
)

type ParseMachine struct {
	state   MachineState
	heading string
}

func (machine *ParseMachine) Exec(arg string, conf *Configuration) error {
	switch machine.state {
	case MachineStateHeading:
		switch machine.heading {
		case "-h":
			conf.host = arg
		case "-u":
			conf.username = arg
		case "-d":
			conf.database = arg
		case "-p":
			conf.passwd = arg
		case "-P":
			conf.port = arg
		default:
			panic(fmt.Sprintf("Invalid option %s", machine.heading))
		}
		machine.state = MachineStateContent
	case MachineStateContent:
		switch arg {
		case "-h", "-u", "-d", "-p", "-c", "-P":
			machine.heading = arg
			machine.state = MachineStateHeading
		// case "-P":
		// 	conf.mode = ConfModeEnterPassword
		case "--help", "help":
			conf.mode = ConfModeUsage
		case "--rails":
			conf.mode = ConfModeRails
		default:
			conf.mode = ConfModeInvalid
		}
	}
	return nil
}

func Config() *Configuration {
	conf := Configuration{}
	conf.host = "localhost"
	conf.port = "3306"
	conf.username = "root"
	conf.passwd = ""
	conf.ParseArgs()
	return &conf
}

func (conf *Configuration) ParseArgs() {
	args := os.Args[1:]
	parseMachine := ParseMachine{}
	for _, arg := range args {
		parseMachine.Exec(arg, conf)
	}
}

func (conf Configuration) Connect() *Window {
	adapter.username = conf.username
	adapter.passwd = conf.passwd
	adapter.host = conf.host
	adapter.port = conf.port

	var dash *DashboardView
	var main_view *ListView
	var window *Window
	width := ui.TermWidth()
	height := ui.TermHeight()

	if conf.database == "" {
		connection = adapter.Connection()
		databases := Databases(connection)
		database_view_list := []ItemView{}
		for _, db := range databases {
			database_view_list = append(database_view_list, ItemView{object: db})
		}

		main_view = NewListView(0, 3, width, height-3, "Select DataBase", database_view_list)
		dash = NewDashboardView(0, 0, width, 3)
		dash.delegate = main_view
		operatios := map[string]string{
			"s":     "Search",
			"c":     "Quick Choose",
			"d":     "Database Detail",
			"C-c":   "Quit",
			"Enter": "Use",
		}

		tips_str := ""
		for key, op := range operatios {
			tips_str += "[" + key + "] " + "[" + op + "]" + "(fg-white,bg-blue)  "
		}
		dash.tips = tips_str
		window = NewWindow(dash, main_view)
	} else {
		db := DatabaseModel{}
		db.name = conf.database
		connection = adapter.Use(db)
		window = db.EntryPoint()
	}
	return window
}
