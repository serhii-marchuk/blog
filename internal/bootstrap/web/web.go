package web

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/serhii-marchuk/blog/internal/ports/web"
	"go.uber.org/fx"
	"html/template"
	"log/slog"
	"net/http"
	"time"
)

var pages = []string{
	"home",
	"about",
	"blog",
	"contact",
	"webhook",
}

var btp = "web/templates/base.html"

func NewWebServer() *echo.Echo {
	return echo.New()
}

func Setup(
	e *echo.Echo,
	r *TemplateRenderer,
	h *web.WebHandler,
) {
	e.HideBanner = true

	for _, pn := range pages {
		r.AddTemplate(pn, template.Must(template.ParseFiles(GetTemplate(pn), btp)))
	}

	e.Renderer = r
	h.Setup(e)
}

func GetTemplate(page string) string {
	return fmt.Sprintf("web/content/%s.html", page)
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
