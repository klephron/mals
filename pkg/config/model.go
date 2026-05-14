package config

type Model struct {
	Name string
	Api  ModelApi
}

type ModelApi interface {
	ModelApiKind() string
}

type ModelApiOpenai struct {
	Url string
}

func (s *ModelApiOpenai) ModelApiKind() string {
	return "openai"
}

func (s *Model) Default() {
}
