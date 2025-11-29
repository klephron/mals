package state

import (
	"github.com/puzpuzpuz/xsync/v4"
	listener "mals/internal/listener/common"
	log "mals/internal/log/common"
)

type State struct {
	EventChan chan Event // should I close this channel?
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
