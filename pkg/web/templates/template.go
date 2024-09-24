package templates

import (
	"html/template"
	"net/http"
)

type Template struct {
	templates *template.Template
}

func NewTemplate() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("static/*.html")),
	}
}

func (t *Template) Render(w http.ResponseWriter, name string, data any) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
