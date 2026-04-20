package config

type Model struct {
	Name string
	Api  ModelApi
}

type ModelApi interface {
	ModelApiKind() string
}

type ModelApiOpenai struct {
	Url         string
	MaxTokens   *int32
	Temperature *float32
}

func (s *ModelApiOpenai) ModelApiKind() string {
	return "openai"
}
