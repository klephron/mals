package state

import (
	"context"
	listener "mals/internal/listener/common"
)

type ListenerValue struct {
	cancel context.CancelFunc
}

func (s *State) ListenerAdd(listener listener.Listener) {
	s.Listeners.Store(listener, nil)
}

func (s *State) ListenerDelete(listener listener.Listener) bool {
	value, loaded := s.Listeners.LoadAndDelete(listener)
	if !loaded {
		return false
	}
	value.cancel()
	return true
}

func (s *State) ListenerListen(listener listener.Listener, ctx context.Context) error {
	lctx, cancel := context.WithCancel(ctx)
	s.Listeners.Store(listener, &ListenerValue{cancel: cancel})

	if err := listener.ListenAndServe(lctx); err != nil {
		return err
	}

	s.EventChan <- &EventListenerDone{}
	return nil
}
