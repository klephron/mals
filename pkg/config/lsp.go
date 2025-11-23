package config

type LspSpecStdio struct {
	LspSpec
	Cmd []string `json:"cmd"`
}
