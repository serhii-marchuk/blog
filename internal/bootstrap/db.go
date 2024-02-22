package bootstrap

import (
	"context"
	"github.com/serhii-marchuk/blog/internal/bootstrap/configs"
	"github.com/serhii-marchuk/blog/internal/bootstrap/web"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"os"
)

func NewDb(l *web.AppLogger) *gorm.DB {
	cfg := configs.NewDbConfig(l)
	db, err := gorm.Open(postgres.Open(cfg.GetDns()), &gorm.Config{})

	if err != nil {
		l.Logger.LogAttrs(context.Background(), slog.LevelError, "Error init DB", slog.String("err", err.Error()))
		os.Exit(0)
	}

	return db
}
