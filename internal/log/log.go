package log

import (
	"fmt"
	"log/slog"
	"mals/pkg/config"
)

type Log struct {
	logFiles []*logFile
}

func getLevel(level string) (slog.Level, error) {
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

func Open(logs []config.Log) (*Log, error) {
	l := &Log{
		logFiles: []*logFile{},
	}

	for _, log := range logs {
		switch logT := log.(type) {
		case *config.LogFile:
			opened, err := openFile(logT.File, logT.Level)

			if err != nil {
				l.Close()
				return nil, err
			}

			l.logFiles = append(l.logFiles, opened)

		default:
			l.Close()
			return nil, fmt.Errorf("unhandled log type %T", logT)
		}
	}

	return l, nil
}

func (l *Log) Close() {
	for _, log := range l.logFiles {
		log.file.Close()
	}
}

func (l *Log) Debug(msg string, args ...any) {
	for _, log := range l.logFiles {
		log.logger.Debug(msg, args...)
	}
}

func (l *Log) Info(msg string, args ...any) {
	for _, log := range l.logFiles {
		log.logger.Info(msg, args...)
	}
}

func (l *Log) Warn(msg string, args ...any) {
	for _, log := range l.logFiles {
		log.logger.Warn(msg, args...)
	}
}

func (l *Log) Error(msg string, args ...any) {
	for _, log := range l.logFiles {
		log.logger.Error(msg, args...)
	}
}
