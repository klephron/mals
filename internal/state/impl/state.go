package state

import (
	listener "mals/internal/listener"
	log "mals/internal/log"

	"github.com/puzpuzpuz/xsync/v4"
)

type StateImpl struct {
	EventChan chan Event
	Listeners *xsync.Map[listener.Listener, ListenerState]
	Logs      *xsync.Map[log.Log, struct{}]
}

func New() *StateImpl {
	return &StateImpl{
		EventChan: make(chan Event),
		Listeners: xsync.NewMap[listener.Listener, ListenerState](),
		Logs:      xsync.NewMap[log.Log, struct{}](),
	}
}

func (s *StateImpl) Wait() {
	s.EventLoop()
}

func (s *StateImpl) Close() {
	s.Listeners.Range(func(key listener.Listener, value ListenerState) bool {
		s.ListenerDelete(key)
		return true
	})
	s.Logs.Range(func(key log.Log, value struct{}) bool {
		s.LogDelete(key)
		return true
	})
	close(s.EventChan)
}
