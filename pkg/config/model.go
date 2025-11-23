package config

type ModelSpecOpenAI struct {
	ModelSpec
	Url         string  `json:"url"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float32 `json:"temperature"`
}
