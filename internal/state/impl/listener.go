package state

import (
	"context"
	listener "mals/internal/listener"
)

func (s *StateImpl) ListenerGet(listener listener.Listener) ListenerValue {
	if value, ok := s.Listeners.Load(listener); ok {
		return value
	}
	return nil
}

func (s *StateImpl) ListenerSet(listener listener.Listener, value ListenerValue) {
	s.Listeners.Store(listener, value)
}

func (s *StateImpl) ListenerAdd(listener listener.Listener) {
	s.Listeners.Store(listener, nil)
}

func (s *StateImpl) ListenerDelete(listener listener.Listener) bool {
	value, loaded := s.Listeners.LoadAndDelete(listener)
	if !loaded {
		return false
	}
	value.Cancel()
	return true
}

func (s *StateImpl) ListenerListen(listener listener.Listener, ctx context.Context) error {
	lctx, _ := context.WithCancel(ctx)

	// s.Listeners.Store(listener, ListenerValueNew(listener, cancel))

	if err := listener.Listen(lctx); err != nil {
		s.EventChan <- &EventListenerDone{}
		return err
	}

	s.EventChan <- &EventListenerDone{}
	return nil
}
