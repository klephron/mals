package state

import (
	listener "mals/internal/listener"
	log "mals/internal/log"

	"github.com/puzpuzpuz/xsync/v4"
)

type StateImpl struct {
	EventChan chan Event
	Listeners *xsync.Map[listener.Listener, StateListener]
	Logs      *xsync.Map[log.Log, struct{}]
}

func New() *StateImpl {
	return &StateImpl{
		EventChan: make(chan Event),
		Listeners: xsync.NewMap[listener.Listener, StateListener](),
		Logs:      xsync.NewMap[log.Log, struct{}](),
	}
}

func (s *StateImpl) Wait() {
	s.Loop()
}

func (s *StateImpl) Close() {
	s.Listeners.Range(func(key listener.Listener, value StateListener) bool {
		s.ListenerDelete(key)
		return true
	})
	s.Logs.Range(func(key log.Log, value struct{}) bool {
		s.LogDelete(key)
		return true
	})
	close(s.EventChan)
}

func (s *StateImpl) listenerAnyListening() bool {
	r := false
	s.Listeners.Range(func(listener listener.Listener, value StateListener) bool {
		if listener.Listening() {
			r = true
			return false
		}
		return true
	})
	return r
}

func (s *StateImpl) Loop() {
	for {
		event, ok := <-s.EventChan
		if !ok {
			return
		}

		switch event := event.(type) {
		case *EventShutdown:
			return

		case *EventListenerListen:
			go func(event *EventListenerListen) {
				listener := event.listener
				listener.Listen(event.ctx)

				if !s.listenerAnyListening() {
					s.EventChan <- &EventShutdown{}
				}
			}(event)

		case *EventClientLspListen:
			go func(event *EventClientLspListen) {
				client := event.client
				client.Listen(event.ctx)
			}(event)
		}
	}
}
