package config

type Model struct {
	Name     string
	Settings ModelSettings
}

type ModelSettings interface {
	Kind() string
}

type ModelSettingsOpenAI struct {
	ModelSettings
	Url         string
	MaxTokens   int
	Temperature float32
}

func (s *ModelSettingsOpenAI) Kind() string {
	return "openai"
}
