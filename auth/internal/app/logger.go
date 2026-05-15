package app

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func NewLogger() *Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)
	return &Logger{
		Logger: logger,
	}
}
