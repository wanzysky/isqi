package main

import (
	"bufio"
	"fmt"
	ui "github.com/gizak/termui"
	"gopkg.in/yaml.v2"
	"image"
	adpt "github.com/wanzysky/isqi/adapters"
	m "github.com/wanzysky/isqi/models"
	v "github.com/wanzysky/isqi/views"
	wd "github.com/wanzysky/isqi/windows"
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

var SupportedAdapter map[string]bool = map[string]bool{
	"mysql2": true,
	"mysql":  true,
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
		case "--help", "help":
			conf.mode = ConfModeUsage
		case "--rails":
			conf.mode = ConfModeRails
		case "--pswd":
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

func (conf *Configuration) Connect() wd.Naviable {
	switch conf.mode {
	case ConfModeNormal:
	case ConfModeUsage:
		conf.Usage()
	case ConfModeRails:
		conf.Rails()
	case ConfModeConfigFile:
		conf.ConfigFile(conf.sourceFilePath, "yaml")
	case ConfModeEnterPassword:
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter you password:")
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err.Error())
		}
		conf.passwd = text
	case ConfModeInvalid:
	}

	adpt.Adpt.Username = conf.username
	adpt.Adpt.Passwd = conf.passwd
	adpt.Adpt.Host = conf.host
	adpt.Adpt.Port = conf.port

	var dash *v.DashboardView
	var main_view *v.ListView
	width := ui.TermWidth()
	height := ui.TermHeight()

	if conf.database == "" {
		adpt.Connection()
		databases := m.Databases()
		database_view_list := []v.ItemView{}
		for _, db := range databases {
			database_view_list = append(database_view_list, v.ItemView{Object: db})
		}

		main_view = v.NewListView(image.Rect(0, 3, width, height), "Select DataBase", database_view_list)
		dash = v.NewDashboardView(image.Rect(0, 0, width, 3))
		dash.Delegate = main_view
		return wd.NewListWindow(main_view, dash)
	} else {
		db := m.DatabaseModel{}
		db.Name = conf.database
		return wd.NewTableIndexWindow(&db)
	}
}

func (conf *Configuration) Usage() {
	fmt.Print(USAGE)
	os.Exit(0)
}

func (conf *Configuration) Rails() {
	path, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	conf.ConfigFile(path+"/config/database.yml", "yaml")
}

func (conf *Configuration) ConfigFile(path string, extname string) {
	file, err := os.Open(path)
	if err != nil {
		panic("Failed to open file at " + path)
	}

	data := make([]byte, 1024)
	count := 0

	count, err = file.Read(data)
	if err != nil {
		panic("Failed to read config file")
	}

	if count >= 1024 {
		panic("Config file can't be larger than 1M.")
	}

	switch extname {
	case "yaml":
		content_map := make(map[interface{}]interface{})
		err = yaml.Unmarshal(data[:count], &content_map)
		if err != nil {
			panic(err.Error())
		}

		environment := "development"
		env := os.Getenv("RAILS_ENV")
		if env != "" {
			environment = env
		}

		if config_map, ok := content_map[environment].(map[interface{}]interface{}); ok {
			conf.InitFromMap(config_map)
		}
	case "json":
	}
}

func (conf *Configuration) InitFromMap(config_map map[interface{}]interface{}) {
	if config_map["host"] != nil {
		if host, ok := config_map["host"].(string); ok {
			conf.host = host
		}
	}

	if config_map["username"] != nil {
		if username, ok := config_map["username"].(string); ok {
			conf.username = username
		}
	}

	if config_map["password"] != nil {
		if password, ok := config_map["password"].(string); ok {
			conf.passwd = password
		}
	}

	if config_map["database"] != nil {
		if database, ok := config_map["database"].(string); ok {
			conf.database = database
		}
	}

	if config_map["adapter"] != nil {
		if adapter, ok := config_map["adapter"].(string); ok {
			if !SupportedAdapter[adapter] {
				panic(fmt.Sprintf("Adapter %s is not supported yet"))
			}
		}
	}
}
