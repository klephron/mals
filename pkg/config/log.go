package config

import "mals/pkg/core"

type Log struct {
	Name   string
	Level  core.LogLevel
	Output LogOutput
}

type LogOutput interface {
	LogOutputKind() string
}

type LogOutputFile struct {
	File string
}

func (s *LogOutputFile) LogOutputKind() string {
	return "file"
}

func (s *Log) Default() {
	if s.Level == "" {
		s.Level = core.LogLevelInfo
	}
	switch co := s.Output.(type) {
	case *LogOutputFile:
		if co.File == "" {
			co.File = "/dev/stdout"
		}
	}
}
