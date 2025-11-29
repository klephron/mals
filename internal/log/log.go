package log

import (
	"fmt"
	"log/slog"
)

type Log interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Close() error
}

func GetLevel(level string) (slog.Level, error) {
	switch level {
	case "error":
		return slog.LevelError, nil
	case "warn":
		return slog.LevelWarn, nil
	case "info":
		return slog.LevelInfo, nil
	case "debug":
		return slog.LevelDebug, nil
	default:
		return slog.LevelDebug, fmt.Errorf("unexpected log level %v", level)
	}
}
