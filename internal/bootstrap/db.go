package bootstrap

import (
	"context"
	"github.com/serhii-marchuk/blog/internal/bootstrap/configs"
	"github.com/serhii-marchuk/blog/internal/bootstrap/web"
	"gorm.io/driver/postgres"
	"log/slog"
	"os"

	"gorm.io/gorm"
)

func NewDb(l *web.AppLogger, cfg *configs.Configs) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.Database.GetDns()), &gorm.Config{})

	if err != nil {
		l.Logger.LogAttrs(context.Background(), slog.LevelError, "Error init DB", slog.String("err", err.Error()))
		os.Exit(0)
	}

	return db
}
