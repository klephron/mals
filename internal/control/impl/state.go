package state

import (
	"context"
	"fmt"
	listener "mals/internal/listener"
	"mals/internal/listener/lsp"
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

func (s *StateImpl) Loop(ctx context.Context) {
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

				lctx, cancel := context.WithCancel(ctx)

				var state StateListener
				switch listener := listener.(type) {
				case *lsp.ListenerLsp:
					state = NewStateListenerLsp(cancel)
				default:
					s.Warn(fmt.Sprintf("unhandled listener type %T when starting listening", listener))
					state = NewStateListener(cancel)
				}

				s.Listeners.Store(listener, state)
				listener.Listen(lctx)

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
