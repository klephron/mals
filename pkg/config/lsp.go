package config

type Lsp struct {
	Name     string
	Settings LspSpec
}

type LspSpec interface {
	lspspec()
}

type LspSettingsStdio struct {
	LspSpec
	Cmd []string
}
