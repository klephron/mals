package config

type Lsp struct {
	Name     string
	Settings LspSettings
}

type LspSettings interface {
	Kind() string
}

type LspSettingsStdio struct {
	LspSettings
	Cmd []string
}

func (s *LspSettingsStdio) Kind() string {
	return "stdio"
}
