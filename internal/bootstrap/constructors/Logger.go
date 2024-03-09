package constructors

import (
	"log/slog"
	"os"
)

type Logger struct {
	Logger *slog.Logger
}

func NewLogger() *Logger {
	return &Logger{Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil))}
}
