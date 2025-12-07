package main

import (
	"context"
	"mals/internal/plane/plane"
)

func main() {
	params := argParse()

	config, err := configLoad(&params)
	if err != nil {
		panic(err)
	}

	ctx, cancel := signalHandle(context.Background())
	defer cancel()

	plane := plane.New()

	plane.Serve(func() {
		go func() {
			<-ctx.Done()
			plane.Lifecycle().Shutdown()
		}()

		configInit(config, plane)
		configLog(config, plane)
	})
}
