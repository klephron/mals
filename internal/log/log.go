package log

import (
	"fmt"
	"log/slog"
	"mals/pkg/config"
	"os"
)

type logFile struct {
	logfile *os.File
	logger  *slog.Logger
}

type Log struct {
	logFiles []*logFile
}

func slogLevel(level string) (slog.Level, error) {
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

func openFile(path string, level string) (*logFile, error) {
	logfile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		return nil, fmt.Errorf("unable to open log file %s", path)
	}

	slogLevel, err := slogLevel(level)
	if err != nil {
		logfile.Close()
		return nil, err
	}

	log := slog.New(slog.NewTextHandler(logfile, &slog.HandlerOptions{
		AddSource: true,
		Level:     slogLevel,
	}))

	return &logFile{logfile: logfile, logger: log}, nil
}

func Open(loggers []config.Log) (*Log, error) {
	l := &Log{logFiles: []*logFile{}}

	for _, logger := range loggers {
		switch typed := logger.(type) {
		case *config.LogFile:
			opened, err := openFile(typed.File, typed.Level)

			if err != nil {
				l.Close()
				return nil, err
			}

			l.logFiles = append(l.logFiles, opened)

		default:
			l.Close()
			return nil, fmt.Errorf("unhandled log type %T", typed)
		}
	}

	return l, nil
}

func (l *Log) Close() {
	for _, logger := range l.logFiles {
		logger.logfile.Close()
	}
}

func (l *Log) Debug(msg string, args ...any) {
	for _, logger := range l.logFiles {
		logger.logger.Debug(msg, args...)
	}
}

func (l *Log) Info(msg string, args ...any) {
	for _, logger := range l.logFiles {
		logger.logger.Info(msg, args...)
	}
}

func (l *Log) Warn(msg string, args ...any) {
	for _, logger := range l.logFiles {
		logger.logger.Warn(msg, args...)
	}
}

func (l *Log) Error(msg string, args ...any) {
	for _, logger := range l.logFiles {
		logger.logger.Error(msg, args...)
	}
}
