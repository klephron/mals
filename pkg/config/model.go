package config

type ModelSpec interface {
	modelspec()
}

type Model struct {
	Name     string
	Settings ModelSpec
}

type ModelSpecOpenAI struct {
	ModelSpec
	Url         string
	MaxTokens   int
	Temperature float32
}
