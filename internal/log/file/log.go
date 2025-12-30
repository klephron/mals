package file

import (
	"fmt"
	"log/slog"
	"mals/internal/log"
	"os"
)

type LogFile struct {
	log.Log
	file   *os.File
	logger *slog.Logger
}

func Open(path string, levelString string) (*LogFile, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		return nil, fmt.Errorf("unable to open log file %s", path)
	}

	level, err := log.GetLevel(levelString)
	if err != nil {
		file.Close()
		return nil, err
	}

	slogLevel := slog.LevelInfo
	switch level {
	case log.LevelDebug:
		slogLevel = slog.LevelDebug
	case log.LevelInfo:
		slogLevel = slog.LevelInfo
	case log.LevelWarn:
		slogLevel = slog.LevelWarn
	case log.LevelError:
		slogLevel = slog.LevelError
	}

	logger := slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
		AddSource: false,
		Level:     slogLevel,
	}))

	return &LogFile{file: file, logger: logger}, nil
}

func (s *LogFile) Debug(msg string, args ...any) {
	s.logger.Debug(msg, args...)
}

func (s *LogFile) Info(msg string, args ...any) {
	s.logger.Info(msg, args...)
}

func (s *LogFile) Warn(msg string, args ...any) {
	s.logger.Warn(msg, args...)
}

func (s *LogFile) Error(msg string, args ...any) {
	s.logger.Error(msg, args...)
}

func (s *LogFile) Close() error {
	return s.file.Close()
}
