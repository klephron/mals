package config

type Lsp struct {
	Name string
	Api  LspApi
}

type LspApi interface {
	LspApiKind() string
}

type LspApiStdio struct {
	Cmd []string
}

func (s *LspApiStdio) LspApiKind() string {
	return "stdio"
}
