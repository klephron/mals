package state

import (
	"context"
	"mals/internal/client"
	"mals/internal/listener"
	"mals/internal/log"
	"mals/pkg/config"

	"github.com/puzpuzpuz/xsync/v4"
)

type State struct {
	Logs      *xsync.Map[string, *LogValue]
	Listeners *xsync.Map[string, *ListenerValue]
	Clients   *xsync.Map[client.Client, *ClientValue]
}

type LogValue struct {
	Config  config.Log
	Log     log.Log
	Enabled bool
}

type ListenerValue struct {
	Config     config.Listener
	Listener   listener.Listener
	CancelFunc context.CancelFunc
	Clients    *xsync.Map[client.Client, struct{}]
}

type ClientValue struct {
	Listener   listener.Listener
	CancelFunc context.CancelFunc
}
