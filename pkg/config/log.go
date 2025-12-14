package config

type Log interface {
	Name() string
	Level() string
}

type LogGeneric struct {
	Log   `json:"log,omitempty"`
	name  string
	level string
}

func (s *LogGeneric) Name() string {
	return s.name
}

func (s *LogGeneric) Level() string {
	return s.level
}

func NewLogGeneric(name string, level string) LogGeneric {
	return LogGeneric{
		name:  name,
		level: level,
	}
}

type LogFile struct {
	LogGeneric
	File string
}
