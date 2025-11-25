package file

import (
	"fmt"
	"log/slog"
	"mals/internal/log/common"
	"os"
)

type LogFile struct {
	common.Log
	file   *os.File
	logger *slog.Logger
}

func Open(path string, level string) (*LogFile, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		return nil, fmt.Errorf("unable to open log file %s", path)
	}

	lvl, err := common.GetLevel(level)
	if err != nil {
		file.Close()
		return nil, err
	}

	logger := slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
		AddSource: true,
		Level:     lvl,
	}))

	return &LogFile{file: file, logger: logger}, nil
}

func (l *LogFile) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *LogFile) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *LogFile) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *LogFile) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *LogFile) Close() error {
	return l.file.Close()
}
