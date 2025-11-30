package main

import (
	"context"
	"os/signal"
	"syscall"
)

func signalHandle(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
}
