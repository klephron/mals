package config

type Lsp struct {
	Name     string
	Settings LspSpec
}

type LspSpec interface {
	lspspec()
}

type LspSpecStdio struct {
	LspSpec
	Cmd []string
}
