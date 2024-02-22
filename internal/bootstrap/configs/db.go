package configs

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/serhii-marchuk/blog/internal/bootstrap/web"
	"log/slog"
)

type DbConfig struct {
	Host          string `env:"DB_HOST"`
	User          string `env:"DB_USER"`
	Password      string `env:"DB_PASSWORD"`
	DBName        string `env:"DB_NAME"`
	Port          int    `env:"DB_PORT"`
	SSLMode       string `env:"DB_SSL_MODE"`
	MigrationPath string `env:"DB_MIGRATIONS_PATH"`
}

func NewDbConfig(l *web.AppLogger) *DbConfig {
	cfg := &DbConfig{}
	if err := env.Parse(cfg); err != nil {
		l.Logger.LogAttrs(context.Background(), slog.LevelError, "Error read .env file", slog.String("err", err.Error()))
	}

	return cfg
}

func (pc *DbConfig) GetDns() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		pc.Host,
		pc.User,
		pc.Password,
		pc.DBName,
		pc.Port,
		pc.SSLMode,
	)
}

func (pc *DbConfig) GetDbSource() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s", // "postgresql://app:app@postgres:5432/app?sslmode=disable"
		pc.User,
		pc.Password,
		pc.Host,
		pc.Port,
		pc.DBName,
		pc.SSLMode,
	)
}
