package constructors

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/serhii-marchuk/blog/internal/bootstrap/configs"
	"log/slog"
	"os"
	"strings"
	"text/template"
	"time"
)

type Renderer struct {
	NavCfg *configs.NavConfig
	Logger *slog.Logger
}

type Page struct {
	NavBar string
	Body   string
}

type Navigation struct {
	Items []NavItem
}

type NavItem struct {
	Name      string
	Uri       string
	ClassName string
}

func NewRenderer(cfg *configs.WebCfg, l *Logger) *Renderer {
	return &Renderer{NavCfg: cfg.NavCfg, Logger: l.Logger}
}

func (r *Renderer) Render(rc *redis.Client, pageName string) ([]byte, error) {
	var html []byte
	var err error

	html, _ = rc.Get(context.Background(), pageName).Bytes()
	if html == nil {
		p := Page{NavBar: r.renderNavBar(pageName), Body: r.renderBody(pageName)}
		html, err = r.renderTemplate(p, r.NavCfg.Layout)
		if err != nil {
			return nil, err
		}
		rdsErr := rc.Set(context.Background(), pageName, html, 24*time.Hour).Err()
		if rdsErr != nil {
			return nil, rdsErr
		}
	}

	return html, nil
}

func (r *Renderer) renderTemplate(data any, tmplPath string) ([]byte, error) {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *Renderer) renderNavBar(page string) string {
	var navigations []NavItem

	for _, item := range r.NavCfg.Pages {
		navigations = append(navigations, NavItem{
			Name:      strings.ToUpper(item),
			Uri:       r.getUri(item),
			ClassName: r.getClassName(item, page),
		})
	}

	nav := Navigation{Items: navigations}
	nb, err := r.renderTemplate(nav, r.NavCfg.NavPath)
	if err != nil {
		r.Logger.LogAttrs(context.Background(), slog.LevelError, "RENDERER_ERROR",
			slog.String("error-msg", err.Error()),
		)
		return ""
	}

	return string(nb)
}

func (r *Renderer) renderBody(page string) string {
	bodyPagePath := r.getTmplPath(page)
	if _, err := os.Stat(bodyPagePath); errors.Is(err, os.ErrNotExist) {
		bodyPagePath = fmt.Sprintf(r.NavCfg.ErrorFilePath, 404)
	}

	body, err := r.renderTemplate(map[string]interface{}{}, bodyPagePath)
	if err != nil {
		r.Logger.LogAttrs(context.Background(), slog.LevelError, "RENDERER_ERROR",
			slog.String("error-msg", err.Error()),
		)
		return ""
	}

	return string(body)
}

func (r *Renderer) getUri(pageName string) string {
	return fmt.Sprintf("/%s", pageName)
}

func (r *Renderer) getClassName(item string, pageName string) string {
	if item == pageName {
		return "nav__item active"
	}
	return "nav__item"
}

func (r *Renderer) getTmplPath(pageName string) string {
	return fmt.Sprintf(r.NavCfg.FilePath, pageName)
}
