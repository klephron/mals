package config

import (
	"fmt"
	"mals/pkg/config"
)

type Model struct {
	Name string    `mapstructure:"name"`
	Api  *ModelApi `mapstructure:"api"`
}

type ModelApi struct {
	Kind ModelApiKind `mapstructure:"kind"`
	Url  *string      `mapstructure:"url"`
}

type ModelApiKind string

const (
	ModelApiKindOpenai ModelApiKind = "openai"
)

func (o *Model) Wire(c *config.Model) error {
	o.Name = c.Name

	switch ca := c.Api.(type) {
	case *config.ModelApiOpenai:
		o.Api = &ModelApi{
			Kind: ModelApiKindOpenai,
		}

		if ca.Url != "" {
			o.Api.Url = &ca.Url
		}

	default:
		return fmt.Errorf("unknown model settings kind")
	}

	return nil
}

func (o *Model) Unwire() (*config.Model, error) {
	c := &config.Model{
		Name: o.Name,
	}

	if o.Api != nil {
		switch o.Api.Kind {
		case ModelApiKindOpenai:
			api := &config.ModelApiOpenai{}
			if o.Api.Url != nil {
				api.Url = *o.Api.Url
			}
			c.Api = api
		default:
			c.Api = nil
		}
	}

	return c, nil
}
