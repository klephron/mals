package config

type Model struct {
	Name string
	Api  ModelApi
}

type ModelApi interface {
	Kind() string
}

type ModelApiOpenai struct {
	Url         string
	MaxTokens   int
	Temperature float32
}

func (s *ModelApiOpenai) Kind() string {
	return "openai"
}
