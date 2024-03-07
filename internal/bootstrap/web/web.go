package web

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/serhii-marchuk/blog/internal/bootstrap/configs"
	"github.com/serhii-marchuk/blog/internal/ports/web"
	"go.uber.org/fx"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func NewWebServer() *echo.Echo {
	return echo.New()
}

func Setup(
	e *echo.Echo,
	h *web.WebHandler,
) {
	e.HideBanner = true
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		LogMethod:   true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	h.Setup(e)
}

func Start(lc fx.Lifecycle, e *echo.Echo, cfg *configs.Configs) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := e.Start(fmt.Sprintf(":%d", cfg.WebPort)); !errors.Is(err, http.ErrServerClosed) {
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
