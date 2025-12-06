package log

import (
	"fmt"
)

type Log interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Close() error
}

type Level int

const (
	LevelAll   Level = iota
	LevelTrace Level = iota
	LevelDebug Level = iota
	LevelInfo  Level = iota
	LevelWarn  Level = iota
	LevelError Level = iota
	LevelOff   Level = iota
)

func GetLevel(level string) (Level, error) {
	switch level {
	case "error":
		return LevelError, nil
	case "warn":
		return LevelWarn, nil
	case "info":
		return LevelInfo, nil
	case "debug":
		return LevelDebug, nil
	default:
		return LevelOff, fmt.Errorf("unexpected log level %v", level)
	}
}
