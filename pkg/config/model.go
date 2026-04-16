package config

type Model struct {
	Name string
	Api  ModelApi
}

type ModelApi interface {
	ModelApi() string
}

type ModelApiOpenai struct {
	Url         string
	MaxTokens   int
	Temperature float32
}

func (s *ModelApiOpenai) ModelApi() string {
	return "openai"
}
