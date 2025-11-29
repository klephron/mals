package state

import (
	"github.com/puzpuzpuz/xsync/v4"
	listener "mals/internal/listener/common"
	log "mals/internal/log/common"
)

type State struct {
	EventChan chan Event
	Listeners *xsync.Map[listener.Listener, *ListenerValue]
	Logs      *xsync.Map[log.Log, struct{}]
}

func New() *State {
	return &State{
		EventChan: make(chan Event),
		Listeners: xsync.NewMap[listener.Listener, *ListenerValue](),
		Logs:      xsync.NewMap[log.Log, struct{}](),
	}
}

func (s *State) Wait() {
	for range s.EventChan {
		c := xsync.NewCounter()

		// check whether any resource is active
		s.Listeners.Range(func(listener listener.Listener, value *ListenerValue) bool {
			if listener.Listening() {
				c.Inc()
			}
			return true
		})

		if c.Value() == 0 {
			return
		}
	}
}

func (s *State) Close() {
	s.Listeners.Range(func(key listener.Listener, value *ListenerValue) bool {
		s.ListenerDelete(key)
		return true
	})
	s.Logs.Range(func(key log.Log, value struct{}) bool {
		s.LogDelete(key)
		return true
	})
	close(s.EventChan)
}
