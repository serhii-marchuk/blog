package bootstrap

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/serhii-marchuk/blog/internal/bootstrap/configs"
	"log/slog"
	"os"
)

type Migrator struct {
	Type string
}

func NewMigrator(d string) *Migrator {
	return &Migrator{Type: d}
}

func (m *Migrator) RunDbMigration(cfg *configs.DbConfig) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mgrt, err := migrate.New(cfg.MigrationPath, cfg.GetDbSource())
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelError, "Error run migrations", slog.String("err", err.Error()))
	}

	if m.Type == "up" {
		if err := mgrt.Up(); err == nil {
			logger.LogAttrs(context.Background(), slog.LevelInfo, "Migrations successfully applied")
		}
	}

	if m.Type == "down" {
		if err := mgrt.Down(); err == nil {
			logger.LogAttrs(context.Background(), slog.LevelInfo, "Migrations successfully roll-back")
		}
	}
}
