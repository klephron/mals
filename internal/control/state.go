package state

import (
	"context"
	"mals/internal/listener"
	"mals/internal/log"
)

type State interface {
	Wait(ctx context.Context)

	ListenerAdd(listener listener.Listener)
	ListenerDelete(listener listener.Listener)
	ListenerServe(listener listener.Listener, ctx context.Context)

	LogAdd(log log.Log)
	LogDelete(log log.Log)
	LogServe(log log.Log, ctx context.Context)

	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}
