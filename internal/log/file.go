package log

import (
	"fmt"
	"log/slog"
	"os"
)

type logFile struct {
	file   *os.File
	logger *slog.Logger
}

func openFile(path string, level string) (*logFile, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		return nil, fmt.Errorf("unable to open log file %s", path)
	}

	lvl, err := getLevel(level)
	if err != nil {
		file.Close()
		return nil, err
	}

	logger := slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
		AddSource: true,
		Level:     lvl,
	}))

	return &logFile{file: file, logger: logger}, nil
}
