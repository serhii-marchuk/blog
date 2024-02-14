package bootstrap

import (
	"context"
	"github.com/serhii-marchuk/blog/internal/bootstrap/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

func NewDb(l *slog.Logger) (*gorm.DB, error) {
	cfg := configs.NewDbConfig(l)
	db, err := gorm.Open(postgres.Open(cfg.GetDns()), &gorm.Config{})

	if err != nil {
		l.LogAttrs(context.Background(), slog.LevelError, "Error init DB", slog.String("err", err.Error()))
		return nil, err
	}

	return db, nil
}
