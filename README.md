# isqi [![Build Status](https://travis-ci.org/wanzysky/isqi.svg?branch=master)](https://travis-ci.org/wanzysky/isqi)
Better database console

![Screenshot](https://github.com/wanzysky/isqi/blob/master/_screenshots/1.gif)

# QuickStart

### Linux
  `\curl -sSL http://dwz.cn/isqi_linux | bash -s stable`

### MacOS
  `\curl -sSL http://dwz.cn/isqi_mac | bash -s stable`

# description

```
Usage:
	isqi -p <PASSWD>

Options:
	-h Database host address, default 'localhost'
	-u Username of database user, default 'root'
	-p Password of database user, default use no password
	-d Database name of your choice, databases will be listed by default
	-c Path to your YAML configuration file, this will overwrite command line options
	-a Adapter of database
	--rails Get configurations from rails app database configuration file
	--help Show usages

Rutime Operations:
  [S]             Search in list
  [C]             Choose by line number from list or table
  [DirectionKeys] Move cursor
  [C-f]           Page down
  [C-b]           Page up
  [D]             Current Line Detail
  [Enter]         Choose currently hilighted line or row
  [Esc]           Back to previous window
  [C-c]           Exit

```

# TODO
  - Multiple Adapters Support

  - SQL Autocomplement

# More
For more informations please visit https://github.com/wanzysky/isqi

Bug reporting & feature request: https://github.com/wanzysky/isqi/issues

Author connection: i@wanzy.me

