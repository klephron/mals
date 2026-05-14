package config

import (
	"fmt"
	"mals/pkg/config"
)

type Lsp struct {
	Name string  `mapstructure:"name"`
	Api  *LspApi `mapstructure:"api"`
}

type LspApi struct {
	Kind LspApiKind `mapstructure:"kind"`
	Cmd  []string   `mapstructure:"cmd"`
}

type LspApiKind string

const (
	LspApiKindStdio LspApiKind = "stdio"
)

func (o *Lsp) Wire(c *config.Lsp) error {
	o.Name = c.Name

	switch ca := c.Api.(type) {
	case *config.LspApiStdio:
		o.Api = &LspApi{
			Kind: LspApiKindStdio,
			Cmd:  ca.Cmd,
		}
	default:
		return fmt.Errorf("unknown lsp settings kind")
	}

	return nil
}

func (o *Lsp) Unwire() (*config.Lsp, error) {
	c := &config.Lsp{
		Name: o.Name,
	}

	switch o.Api.Kind {
	case LspApiKindStdio:
		c.Api = &config.LspApiStdio{
			Cmd: o.Api.Cmd,
		}
	default:
		c.Api = nil
	}

	return c, nil
}
