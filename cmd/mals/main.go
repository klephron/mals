package main

import (
	"context"
	"mals/internal/plane"
)

func main() {
	params := argParse()

	_, err := configLoad(&params)
	if err != nil {
		panic(err)
	}

	ctx, cancel := signalHandle(context.Background())
	defer cancel()

	plane := plane.New()

	go func() {
		<-ctx.Done()
		plane.Lifecycle.Shutdown()
	}()

	plane.Serve()
}
