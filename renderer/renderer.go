package renderer

import (
	"fmt"
	"html/template"
)

const (
    templatesDir = "views"
)

var Tmpl *template.Template

func Init() {
    var err error
    Tmpl, err = template.ParseGlob(templatesDir+"/*.html")
    if err != nil {
        fmt.Printf("couldn't parse templates: %v\n", err.Error())
    }
} 
