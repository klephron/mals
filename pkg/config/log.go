package config

type Log struct {
	Name  string
	Level string
	Kind  LogKind
}

type LogKind interface {
	Kind() string
}

type LogKindFile struct {
	LogKind
	File string
}

func (s *LogKindFile) Kind() string {
	return "file"
}
