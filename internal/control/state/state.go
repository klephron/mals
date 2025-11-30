package state

import (
	"mals/internal/listener"
	"mals/internal/log"

	"github.com/puzpuzpuz/xsync/v4"
)

type State struct {
	Listeners *xsync.Map[listener.Listener, struct{}]
	Logs      *xsync.Map[log.Log, struct{}]
}

func New() *State {
	return &State{
		Listeners: &xsync.Map[listener.Listener, struct{}]{},
		Logs:      &xsync.Map[log.Log, struct{}]{},
	}
}
