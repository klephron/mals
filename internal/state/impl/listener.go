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

func (s *StateImpl) ListenerFind(listener listener.Listener) (StateListener, error) {
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

	var state StateListener

	switch listener := listener.(type) {
	case *lsp.ListenerLsp:
		state = NewStateListenerLsp(cancel)
	default:
		s.Warn(fmt.Sprintf("unhandled listener type %T when starting listening", listener))
		state = NewStateListener(cancel)
	}

	s.Listeners.Store(listener, state)

	s.EventChan <- &EventListenerListen{listener: listener, ctx: lctx}

	return nil
}

func (s *StateImpl) ListenerAddConn(listener listener.Listener, ctx context.Context, conn net.Conn) error {
	state, err := s.ListenerFind(listener)
	if err != nil {
		return err
	}
	if state == nil {
		s.Error("listener %v state is nil", listener)
		return nil
	}

	switch state := state.(type) {
	case *StateListenerLsp:
		client := client.New(s, conn)

		ctx, cancel := context.WithCancel(ctx)

		clientState := NewClientState(cancel)
		state.Clients.Store(client, clientState)

		s.EventChan <- &EventClientLspListen{client: client, ctx: ctx}

	default:
		s.Warn("ListenerAddConn unhandled listener %v of state type %T", listener, state)
	}

	return nil
}
