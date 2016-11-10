package main

import (
	"os"
	strings "strings"
)

var short_names map[string]string = map[string]string{
	"a": "adapter",
	"p": "port",
	"P": "passwd",
	"u": "user",
	"h": "host",
	"d": "database",
	"c": "certification",
	"H": "help",
	"v": "version",
}

func params() {
	config := make(map[string]string)
	config["adapter"] = "mysql"
	config["username"] = "root"
	config["host"] = "localhost"
	config["port"] = "3306"
	config["passwd"] = ""
	args := os.Args[1:]

	key := ""
	for _, val := range args {
		if len(key) > 0 {
			assign(key, val, config)
			key = ""
		} else {
			key = parse_key(val)
		}
	}
}

func assign(key, value string, config map[string]string) {
	permitting_keys := map[string]bool{
		"adapter":       true,
		"port":          true,
		"passwd":        true,
		"user":          true,
		"host":          true,
		"database":      true,
		"certification": true,
		"help":          true,
		"version":       true,
	}
	if _, ok := permitting_keys[key]; !ok {
		return
	}
	config[key] = value
}

func parse_key(key string) string {
	if strings.HasPrefix(key, "--") {
		return key[2:]
	} else {
		if !strings.HasPrefix(key, "-") {
			return ""
		}
		return short_names[key[1:]]
	}
}
