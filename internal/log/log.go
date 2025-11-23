package log

import (
	"fmt"
	"log/slog"
	"mals/pkg/config"
	"os"
)

type logFile struct {
	Logfile *os.File
	Logger  *slog.Logger
}

type Log struct {
	loggers []*logFile
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

	return &logFile{Logfile: logfile, Logger: log}, nil
}

func Open(loggers []config.Log) (*Log, error) {
	log := &Log{loggers: []*logFile{}}

	for _, logger := range loggers {
		switch typed := logger.(type) {
		case *config.LogFile:
			opened, err := openFile(typed.File, typed.Level)

			if err != nil {
				log.Close()
				return nil, err
			}

			log.loggers = append(log.loggers, opened)

		default:
			log.Close()
			return nil, fmt.Errorf("unhandled log type %T", typed)
		}
	}

	return log, nil
}

func (log *Log) Close() {
	for _, logger := range log.loggers {
		logger.Logfile.Close()
	}
}

func (log *Log) Debug(msg string, args ...any) {
	for _, logger := range log.loggers {
		logger.Logger.Debug(msg, args...)
	}
}

func (log *Log) Info(msg string, args ...any) {
	for _, logger := range log.loggers {
		logger.Logger.Info(msg, args...)
	}
}

func (log *Log) Warn(msg string, args ...any) {
	for _, logger := range log.loggers {
		logger.Logger.Warn(msg, args...)
	}
}

func (log *Log) Error(msg string, args ...any) {
	for _, logger := range log.loggers {
		logger.Logger.Error(msg, args...)
	}
}
