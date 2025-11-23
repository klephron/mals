package log

import (
	"fmt"
	"log/slog"
	"os"
)

type Log struct {
	logfile *os.File
	Logger  *slog.Logger
}

func logLevel(level int) slog.Level {
	switch {
	case level <= 0:
		return slog.LevelError
	case level <= 1:
		return slog.LevelWarn
	case level <= 2:
		return slog.LevelInfo
	default:
		return slog.LevelDebug
	}
}

func Open(path string, level int) (*Log, error) {
	logfile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		return nil, fmt.Errorf("unable to open log file %s", path)
	}

	logger := slog.New(slog.NewTextHandler(logfile, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel(level),
	}))

	return &Log{logfile: logfile, Logger: logger}, nil
}

func (log *Log) Close() {
	log.logfile.Close()
}
