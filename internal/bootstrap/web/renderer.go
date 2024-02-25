package web

import (
	"github.com/labstack/echo/v4"
	"github.com/serhii-marchuk/blog/internal/bootstrap/configs"
	"html/template"
	"io"
)

type TemplateRenderer struct {
	templates map[string]*template.Template
	webCfg    *configs.WebCfg
}

func NewRenderer(cfg *configs.WebCfg) *TemplateRenderer {
	return &TemplateRenderer{templates: make(map[string]*template.Template), webCfg: cfg}
}

func (t *TemplateRenderer) AddTemplate(key string, tmpl *template.Template) {
	t.templates[key] = tmpl
}

func (t *TemplateRenderer) Render(writer io.Writer, s string, i interface{}, ctx echo.Context) error {
	tmpl, ok := t.templates[s]
	if !ok {
		tmpl = template.Must(template.ParseFiles(t.webCfg.NavCfg.GetErrorPageFile("404"), t.webCfg.NavCfg.BaseTemplatePath))
	}

	return tmpl.ExecuteTemplate(writer, "base", i)
}
