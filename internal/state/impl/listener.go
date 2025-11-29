package state

import (
	"context"
	"fmt"
	listener "mals/internal/listener"
	"mals/internal/listener/lsp"
)

func (s *StateImpl) ListenerAdd(listener listener.Listener) {
	s.Listeners.Store(listener, nil)
}

func (s *StateImpl) ListenerDelete(listener listener.Listener) bool {
	value, loaded := s.Listeners.LoadAndDelete(listener)
	if !loaded || value == nil {
		return false
	}
	value.Cancel()
	return true
}

func (s *StateImpl) ListenerListen(listener listener.Listener, ctx context.Context) error {
	if _, found := s.Listeners.Load(listener); !found {
		err := fmt.Errorf("listener %v not found", listener)
		s.Warn(err.Error())
		return err
	}

	lctx, cancel := context.WithCancel(ctx)

	var value ListenerValue

	switch listener.(type) {
	case *lsp.ListenerLsp:
		value = NewListenerValueLsp(cancel)
	default:
		s.Warn(fmt.Sprintf("unhandled listener type %T when starting listening", listener))
		value = NewListenerValueGeneric(cancel)
	}

	s.Listeners.Store(listener, value)

	if err := listener.Listen(lctx); err != nil {
		s.EventChan <- &EventListenerDone{}
		return err
	}

	s.EventChan <- &EventListenerDone{}
	return nil
}
