package main

import (
	"encoding/json"
	"fmt"
	"mals/pkg/config"
	"os"
)

func loadConfig(params *Params) (*config.Config, error) {
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
