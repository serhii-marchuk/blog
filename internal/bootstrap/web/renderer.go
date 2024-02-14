package web

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

type TemplateRenderer struct {
	templates map[string]*template.Template
}

func NewRenderer() *TemplateRenderer {
	return &TemplateRenderer{templates: make(map[string]*template.Template)}
}

func (t *TemplateRenderer) AddTemplate(key string, tmpl *template.Template) {
	t.templates[key] = tmpl
}

func (t *TemplateRenderer) Render(writer io.Writer, s string, i interface{}, ctx echo.Context) error {
	tmpl, ok := t.templates[s]
	if !ok {
		tmpl = template.Must(template.ParseFiles("web/content/404.html", "web/templates/base.html"))
	}

	return tmpl.ExecuteTemplate(writer, "base", i)
}
