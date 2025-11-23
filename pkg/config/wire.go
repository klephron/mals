package config

import (
	"encoding/json"
)

func (c *Config) UnmarshallJSON(data []byte) error {
	init := Config{
		Models: []*Model{},
		Lsps:   []*Lsp{},
		Usages: []*Usage{},
	}
	err := json.Unmarshal(data, &init)
	if err != nil {
		return err
	}
	*c = init
	return nil
}
