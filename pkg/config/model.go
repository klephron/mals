package config

type ModelSpecOpenAI struct {
	ModelSpec
	Url         string
	MaxTokens   int
	Temperature float32
}
