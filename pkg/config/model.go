package config

type ModelSettings interface {
	modelspec()
}

type Model struct {
	Name     string
	Settings ModelSettings
}

type ModelSettingsOpenAI struct {
	ModelSettings
	Url         string
	MaxTokens   int
	Temperature float32
}
