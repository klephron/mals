package main

import (
	"encoding/json"
	"mals/internal/plane"
	"mals/pkg/config"
	"os"
)

func configLoad(params *Params) (*config.Config, error) {
	bytes, err := os.ReadFile(params.Config)
	if err != nil {
		return nil, err
	}

	var c config.Config
	if err := json.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func configInit(config *config.Config, plane plane.Plane) {
	for _, log := range config.Loggers {
		if err := plane.Log().LogRegister(log.Name(), log); err != nil {
			plane.Log().Errorf("%v", err)
		}
		if err := plane.Log().LogCreate(log.Name()); err != nil {
			plane.Log().Errorf("%v", err)
		}
		if err := plane.Log().LogStart(log.Name()); err != nil {
			plane.Log().Errorf("%v", err)
		}
	}

	for _, usage := range config.Usages {
		if usage == nil {
			continue
		}
		if err := plane.Usage().UsageRegister(*usage); err != nil {
			plane.Log().Errorf("%v", err)
		}
	}

	for _, model := range config.Models {
		if model == nil {
			continue
		}
		if err := plane.Scope().ModelRegister(*model); err != nil {
			plane.Log().Errorf("%v", err)
		}
	}

	for _, listener := range config.Listeners {
		if err := plane.Listener().ListenerRegister(listener.Name(), listener); err != nil {
			plane.Log().Errorf("%v", err)
		}
		if err := plane.Listener().ListenerCreate(listener.Name()); err != nil {
			plane.Log().Errorf("%v", err)
		}
		if err := plane.Listener().ListenerStart(listener.Name()); err != nil {
			plane.Log().Errorf("%v", err)
		}
	}

}

func configLog(config *config.Config, plane plane.Plane) {
	bytes, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	plane.Log().Debugf("config: %v", string(bytes))
}
