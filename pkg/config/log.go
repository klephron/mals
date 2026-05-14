package config

type Log struct {
	Name   string
	Level  LogLevel
	Output LogOutput
}

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

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
		s.Level = LogLevelInfo
	}
	switch co := s.Output.(type) {
	case *LogOutputFile:
		if co.File == "" {
			co.File = "/dev/stdout"
		}
	}
}
