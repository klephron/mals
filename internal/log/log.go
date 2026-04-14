package log

import (
	"fmt"
	"mals/pkg/config"
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

func GetLevel(level config.LogLevel) (Level, error) {
	switch level {
	case config.LogLevelError:
		return LevelError, nil
	case config.LogLevelWarn:
		return LevelWarn, nil
	case config.LogLevelInfo:
		return LevelInfo, nil
	case config.LogLevelDebug:
		return LevelDebug, nil
	default:
		return LevelOff, fmt.Errorf("unexpected log level %v", level)
	}
}
