package main

import (
	"encoding/json"
	"fmt"
	"mals/internal/control/controller"
	listener "mals/internal/listener/factory"
	log "mals/internal/log/factory"
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

func configInitLogs(config *config.Config, controller *controller.Controller) {
	for _, loggerConfig := range config.Loggers {
		log, err := log.OpenConfig(loggerConfig)
		if err != nil {
			panic(err)
		}
		controller.LogAdd(log)
	}
}

func configInitListeners(config *config.Config, controller *controller.Controller) {
	for _, listenerConfig := range config.Listeners {
		listener, err := listener.NewConfig(controller, listenerConfig)
		if err != nil {
			panic(err)
		}
		controller.ListenerAdd(listener)
	}
}

func configLog(config *config.Config, controller *controller.Controller) {
	configJson, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	controller.Debug(fmt.Sprintf("config: %v", string(configJson)))
}
