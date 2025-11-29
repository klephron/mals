package state

import (
	"context"
	"fmt"
	listener "mals/internal/listener"
	"mals/internal/listener/lsp"
	client "mals/internal/lsp/client"
	"net"
)

func (s *StateImpl) ListenerAdd(listener listener.Listener) {
	s.Listeners.Store(listener, nil)
}

func (s *StateImpl) ListenerDelete(listener listener.Listener) bool {
	state, loaded := s.Listeners.LoadAndDelete(listener)
	if !loaded || state == nil {
		return false
	}
	state.Cancel()
	return true
}

func (s *StateImpl) ListenerFind(listener listener.Listener) (ListenerState, error) {
	if state, found := s.Listeners.Load(listener); !found {
		err := fmt.Errorf("listener %v not found", listener)
		s.Warn(err.Error())
		return nil, err
	} else {
		return state, nil
	}
}

func (s *StateImpl) ListenerListen(listener listener.Listener, ctx context.Context) error {
	if _, err := s.ListenerFind(listener); err != nil {
		return err
	}

	lctx, cancel := context.WithCancel(ctx)

	var state ListenerState

	switch listener.(type) {
	case *lsp.ListenerLsp:
		state = NewListenerStateLsp(cancel)
	default:
		s.Warn(fmt.Sprintf("unhandled listener type %T when starting listening", listener))
		state = NewListenerStateGeneric(cancel)
	}

	s.Listeners.Store(listener, state)

	if err := listener.Listen(lctx); err != nil {
		s.EventChan <- &EventListenerDown{}
		return err
	}

	s.EventChan <- &EventListenerDown{}
	return nil
}

func (s *StateImpl) ListenerAddConn(listener listener.Listener, ctx context.Context, conn net.Conn) error {
	state, err := s.ListenerFind(listener)
	if err != nil {
		s.Error("listener %v not found", listener)
		return err
	}
	if state == nil {
		s.Error("listener %v state is nil", listener)
		return nil
	}

	switch stateT := state.(type) {
	case *ListenerStateLsp:
		client := client.New(s, conn)

		ctx, cancel := context.WithCancel(ctx)
		state := NewClientState(cancel)

		stateT.Clients.Store(client, state)
		s.EventChan <- &EventClientListen{client: client, ctx: ctx}

	default:
		s.Warn("ListenerAddConn unhandled listener %v of state type %T", listener, stateT)
	}

	return nil
}
