package state

import (
	"context"
	"mals/internal/listener"
	"mals/internal/log"
	"net"
)

type State interface {
	Wait()
	Close()

	ListenerAdd(listener listener.Listener)
	ListenerDelete(listener listener.Listener) bool
	ListenerListen(listener listener.Listener, ctx context.Context) error
	ListenerAddConn(listener listener.Listener, ctx context.Context, conn net.Conn) error

	LogAdd(log log.Log)
	LogDelete(log log.Log) bool
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}
