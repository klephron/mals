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

func (s *Lsp) Default() {
	switch settings := s.Api.(type) {
	case *LspApiStdio:
		if settings.Cmd == nil {
			settings.Cmd = make([]string, 0)
		}
	}
}
