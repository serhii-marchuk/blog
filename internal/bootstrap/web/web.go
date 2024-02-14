package web

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/serhii-marchuk/blog/internal/ports/web"
	"go.uber.org/fx"
	"html/template"
	"log/slog"
	"net/http"
	"time"
)

func NewWebServer() *echo.Echo {
	return echo.New()
}

func Setup(
	e *echo.Echo,
	r *TemplateRenderer,
	h *web.WebHandler,
) {
	e.HideBanner = true

	r.AddTemplate("home", template.Must(template.ParseFiles("web/content/home.html", "web/templates/base.html")))
	r.AddTemplate("about", template.Must(template.ParseFiles("web/content/about.html", "web/templates/base.html")))
	r.AddTemplate("blog", template.Must(template.ParseFiles("web/content/blog.html", "web/templates/base.html")))
	r.AddTemplate("contact", template.Must(template.ParseFiles("web/content/contact.html", "web/templates/base.html")))

	e.Renderer = r
	h.Setup(e)
}

func Start(lc fx.Lifecycle, e *echo.Echo) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := e.Start(":80"); !errors.Is(err, http.ErrServerClosed) {
					slog.Error(err.Error())
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
			defer cancel()

			return e.Shutdown(ctx)
		},
	})
}
