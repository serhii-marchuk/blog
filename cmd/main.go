package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/serhii-marchuk/blog/internal/bootstrap/web"
	webHandl "github.com/serhii-marchuk/blog/internal/ports/web"
	"go.uber.org/fx"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if godotenv.Load("./configs/.env") != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Error read .env file")
	}
	slog.SetDefault(logger)

	//db, err := bootstrap.NewDb(logger)
	//if err != nil {
	//	os.Exit(0)
	//}

	fx.New(
		fx.Provide(web.NewWebServer),
		fx.Provide(web.NewRenderer),
		fx.Provide(webHandl.NewWebHandler),
		fx.Invoke(
			web.Setup,
			web.Start,
		),
	).Run()
}
