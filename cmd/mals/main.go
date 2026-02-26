package main

import (
	"context"
	"mals/internal/plane/planeimpl"
)

func main() {
	params := argParse()

	config, err := configLoad(&params)
	if err != nil {
		panic(err)
	}

	ctx, cancel := signalHandle(context.Background())
	defer cancel()

	plane := planeimpl.New()

	plane.Run(func() {
		go func() {
			<-ctx.Done()
			if err := plane.Shutdown(); err != nil {
				panic(err)
			}
		}()

		configInit(config, plane)
		configLog(config, plane)
	})
}
