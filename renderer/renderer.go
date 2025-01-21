package renderer

import (
	"fmt"
	"html/template"
    "strings"
    "io"
)

const (
    templatesDir = "views"
)

type PageData struct {
    Title string
    Data  any
}

func NewPageData(title string, data any) PageData {
    return PageData{
        Title: title,
        Data: data,
    }
}

func RenderTemplate(writer io.Writer, name string, data any) {
    t := loadTemplate(name)
    t.Execute(writer, PageData{
        Title: name,
        Data: data,
    })
}

func loadTemplate(name string) *template.Template {
    base := fmt.Sprintf("%s/base.html", templatesDir)
    name = fmt.Sprintf("%s/%s.html", templatesDir, name)

    t, err := template.ParseFiles(base, name)
    if err != nil {
        fmt.Printf("couldn't load template: %v\n%v\n", name, err.Error())
    }

    return t
}

func printTemplates(t *template.Template) {
    var sb strings.Builder
    for i, tmp := range t.Templates() {
        fmt.Fprintf(&sb, "tmpl %v: %+v\n", i, tmp.Name())
    }

    fmt.Printf("%v\n", sb.String())
}
