package lsp

import (
	"mals/pkg/config"
)

type SettingsIO struct {
	config.LspSpec
	Cmd []string `json:"cmd" default:"[]"`
}
