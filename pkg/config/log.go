package config

type Log interface {
	Name() string
}

type LogGeneric struct {
	Log  `json:"log,omitempty"`
	name string
}

func (s *LogGeneric) Name() string {
	return s.name
}

func NewLogGeneric(name string) LogGeneric {
	return LogGeneric{
		name: name,
	}
}

type LogFile struct {
	LogGeneric
	Level string
	File  string
}
