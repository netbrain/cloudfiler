package web

import (
	"html/template"
	"strings"
)

var templateFunctions = template.FuncMap{
	"stringsJoin": strings.Join,
}
