package state

import (
	"context"
	"mals/internal/listener"
	"mals/internal/log"

	"github.com/puzpuzpuz/xsync/v4"
)

type State struct {
	Listeners *xsync.Map[listener.Listener, *ListenerValue]
	Logs      *xsync.Map[log.Log, *LogValue]
}

type LogValue struct {
	Enabled bool
}

type ListenerValue struct {
	CancelFunc context.CancelFunc
}
