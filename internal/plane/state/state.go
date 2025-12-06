package state

import (
	"context"
	"mals/internal/listener"
	"mals/internal/log"

	"github.com/puzpuzpuz/xsync/v4"
)

type State struct {
	Logs      *xsync.Map[string, *LogValue]
	Listeners *xsync.Map[string, *ListenerValue]
}

type LogValue struct {
	Log     log.Log
	Enabled bool
}

type ListenerValue struct {
	Listener   listener.Listener
	CancelFunc context.CancelFunc
}
