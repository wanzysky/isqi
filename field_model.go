package main

const FIELD_MODEL_ATTR_COUNT = 9

type FieldModel struct {
	field      string
	types      string
	collation  string
	null       string
	key        string
	defaults   string
	extra      string
	privileges string
	comment    string
}
