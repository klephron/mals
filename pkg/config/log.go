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
	LogOutput() string
}

type LogOutputFile struct {
	File string
}

func (s *LogOutputFile) LogOutput() string {
	return "file"
}
