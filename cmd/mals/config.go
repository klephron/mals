package main

import (
	"encoding/json"
	"fmt"
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

	// intercept to set defaults
	for _, logger := range c.Loggers {
		switch typed := logger.(type) {
		case *config.LogFile:
			if typed.Level == "" {
				typed.Level = "debug"
			}
		default:
			return nil, fmt.Errorf("unhandled log type %T", typed)
		}
	}

	return &c, nil
}

func configInit(config *config.Config, plane *plane.Plane) {
	for _, log := range config.Loggers {
		if err := plane.Log.Register(log); err != nil {
			plane.Log.Errorf("%v", err)
		}
		if err := plane.Log.Create(log.Name()); err != nil {
			plane.Log.Errorf("%v", err)
		}
		if err := plane.Log.Start(log.Name()); err != nil {
			plane.Log.Errorf("%v", err)
		}
	}
}

func configLog(config *config.Config, plane *plane.Plane) {
	bytes, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	plane.Log.Debugf("config: %v", string(bytes))
}
