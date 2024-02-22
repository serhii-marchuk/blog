package web

import (
	"log/slog"
	"os"
)

type AppLogger struct {
	Logger *slog.Logger
}

func NewAppLogger() *AppLogger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	return &AppLogger{Logger: logger}
}
