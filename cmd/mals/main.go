package main

import (
	"context"
	"mals/internal/control/controller"
	"mals/internal/control/event"
	"mals/internal/control/scheduler"
	"mals/internal/control/state"
)

func main() {
	params := argParse()

	config, err := configLoad(&params)
	if err != nil {
		panic(err)
	}

	ctx, cancel := signalHandle(context.Background())
	defer cancel()

	state := state.New()

	bus := event.NewEventBus()
	defer bus.Close()

	scheduler := scheduler.New(state, bus)
	scheduler.Subscribe()

	go func(ctx context.Context) {
		controller := controller.New(state, bus)

		configInitLogs(config, controller)
		configLog(config, controller)

		configInitListeners(config, controller)

		<-ctx.Done()
		controller.Shutdown()
	}(ctx)

	scheduler.ServeSubscribed()
}
