package config

type Lsp struct {
	Name string
	Api  LspApi
}

type LspApi interface {
	Kind() string
}

type LspApiStdio struct {
	Cmd []string
}

func (s *LspApiStdio) Kind() string {
	return "stdio"
}
