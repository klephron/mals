package main

import (
	"context"
	"os/signal"
	"syscall"
)

func signalHandle() (context.Context, context.CancelFunc) {
	return signalHandleContext(context.Background())
}

func signalHandleContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
}
