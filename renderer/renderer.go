package renderer

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"strings"
)

const (
	templatesDir = "views"
	defaultTitle = "gotodo"
	baseTemplate = "base"
	noContentTemplate = "nocontent"
)

type PageData struct {
	Title string
	Data  any
}

func RenderTemplate(writer io.Writer, name string, data *PageData) error {
	t := loadTemplate(name)
	if t == nil {
		return errors.New("No templates to render")
	}

	if data == nil {
		data = &PageData{
			Title: defaultTitle,
			Data:  nil,
		}
	}

	return t.Execute(writer, data)
}

func loadTemplate(name string) *template.Template {
	base := fmt.Sprintf("%s/%s.html", templatesDir, baseTemplate)
	name = fmt.Sprintf("%s/%s.html", templatesDir, name)

	b, err := template.ParseFiles(base)
	if err != nil {
		return nil
	}

	t, err := b.ParseFiles(name)
	if err == nil {
		return t
	}

	name = fmt.Sprintf("%s/%s.html", templatesDir, noContentTemplate)
	t, err = b.ParseFiles(name)
	if err != nil {
		return nil
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
