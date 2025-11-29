package main

import (
	"encoding/json"
	"fmt"
	"mals/internal/listener/factory"
	"mals/internal/log/factory"
	"mals/internal/state/impl"
)

func main() {
	ctx, stop := signalHandle()
	defer stop()

	params := argParse()

	config, err := loadConfig(&params)
	if err != nil {
		panic(err)
	}

	state := state.New()
	defer state.Close()

	for _, loggerConfig := range config.Loggers {
		log, err := log.OpenConfig(loggerConfig)
		if err != nil {
			panic(err)
		}
		state.LogAdd(log)
	}

	configJson, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	state.Debug(fmt.Sprintf("config: %v", string(configJson)))

	for _, listenerConfig := range config.Listeners {
		listener, err := listener.NewConfig(state, listenerConfig)
		if err != nil {
			panic(err)
		}
		state.ListenerAdd(listener)
		go func() {
			state.ListenerListen(listener, ctx)
		}()
	}

	state.Wait()
}
