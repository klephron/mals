package model

import (
	"mals/pkg/config"
)

type SettingsOpenAI struct {
	config.ModelSpec
	Url         *string  `json:"url"`
	MaxTokens   *int     `json:"max_tokens"`
	Temperature *float32 `json:"temperature"`
}
