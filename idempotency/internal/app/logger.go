package app

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func NewLogger(logLevel string) *Logger {
	var lvl slog.Level
	switch logLevel {
	case "debug":
		lvl = slog.LevelDebug
	case "error":
		lvl = slog.LevelError
	case "warn":
		lvl = slog.LevelWarn
	default:
		lvl = slog.LevelInfo

	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     lvl,
		AddSource: true,
	})
	slog := slog.New(handler)

	return &Logger{
		slog,
	}
}
