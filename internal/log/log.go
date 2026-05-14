package log

import (
	"fmt"
	"mals/pkg/core"
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

func GetLevel(level core.LogLevel) (Level, error) {
	switch level {
	case core.LogLevelError:
		return LevelError, nil
	case core.LogLevelWarn:
		return LevelWarn, nil
	case core.LogLevelInfo:
		return LevelInfo, nil
	case core.LogLevelDebug:
		return LevelDebug, nil
	default:
		return LevelOff, fmt.Errorf("unexpected log level %v", level)
	}
}
