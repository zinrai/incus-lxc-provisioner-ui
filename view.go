package main

import (
	_ "embed"
	"html/template"
)

//go:embed page.html
var pageHTML string

var tmpl = template.Must(template.New("page").Parse(pageHTML))
