package config

type Lsp struct {
	Name string
	Api  LspApi
}

type LspApi interface {
	LspApi() string
}

type LspApiStdio struct {
	Cmd []string
}

func (s *LspApiStdio) LspApi() string {
	return "stdio"
}
