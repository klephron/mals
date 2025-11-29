package state

import (
	listener "mals/internal/listener"
	log "mals/internal/log"

	"github.com/puzpuzpuz/xsync/v4"
)

type StateImpl struct {
	EventChan chan Event
	Listeners *xsync.Map[listener.Listener, ListenerValue]
	Logs      *xsync.Map[log.Log, struct{}]
}

func New() *StateImpl {
	return &StateImpl{
		EventChan: make(chan Event),
		Listeners: xsync.NewMap[listener.Listener, ListenerValue](),
		Logs:      xsync.NewMap[log.Log, struct{}](),
	}
}

func (s *StateImpl) Wait() {
	for range s.EventChan {
		c := xsync.NewCounter()

		// check whether any resource is active
		s.Listeners.Range(func(listener listener.Listener, value ListenerValue) bool {
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

func (s *StateImpl) Close() {
	s.Listeners.Range(func(key listener.Listener, value ListenerValue) bool {
		s.ListenerDelete(key)
		return true
	})
	s.Logs.Range(func(key log.Log, value struct{}) bool {
		s.LogDelete(key)
		return true
	})
	close(s.EventChan)
}
