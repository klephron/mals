package state

import (
	listener "mals/internal/listener/common"
	log "mals/internal/log/common"
)

type State struct {
	Listeners []listener.Listener
	Logs      []log.Log
}

func New() *State {
	return &State{
		Listeners: []listener.Listener{},
		Logs:      []log.Log{},
	}
}
