package main

import (
	"encoding/json"
	"mals/internal/plane"
	"mals/pkg/config"
	"mals/pkg/wire"
	"os"
)

func configLoad(params *Params) (*config.Config, error) {
	bytes, err := os.ReadFile(params.Config)
	if err != nil {
		return nil, err
	}

	var wire wire.Config
	if err := json.Unmarshal(bytes, &wire); err != nil {
		return nil, err
	}

	cfg, err := wire.Unwire()
	if err != nil {
		return nil, err
	}

	config.DefaultConfig(cfg)

	return cfg, nil
}

func configInit(config *config.Config, plane plane.Plane) {
	for _, log := range config.Logs {
		if log == nil {
			continue
		}
		if err := plane.Log().Register(log.Name, log); err != nil {
			plane.Errorf("%v", err)
		}
		if err := plane.Log().Create(log.Name); err != nil {
			plane.Errorf("%v", err)
		}
		if err := plane.Log().Start(log.Name); err != nil {
			plane.Errorf("%v", err)
		}
	}

	for _, lsp := range config.Lsps {
		if lsp == nil {
			continue
		}
		if err := plane.Scope().LspRegister(lsp); err != nil {
			plane.Errorf("%v", err)
		}
	}

	for _, model := range config.Models {
		if model == nil {
			continue
		}
		if err := plane.Scope().ModelRegister(model); err != nil {
			plane.Errorf("%v", err)
		}
	}

	for _, usage := range config.Usages {
		if usage == nil {
			continue
		}
		if err := plane.Usage().Register(usage); err != nil {
			plane.Errorf("%v", err)
		}
	}

	for _, listener := range config.Listeners {
		if listener == nil {
			continue
		}
		if err := plane.Listener().Register(listener.Name, listener); err != nil {
			plane.Errorf("%v", err)
		}
		if err := plane.Listener().Create(listener.Name); err != nil {
			plane.Errorf("%v", err)
		}
		if err := plane.Listener().Start(listener.Name); err != nil {
			plane.Errorf("%v", err)
		}
	}

}

func configLog(config *config.Config, plane plane.Plane) {
	bytes, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	plane.Debugf("config: %v", string(bytes))
}
