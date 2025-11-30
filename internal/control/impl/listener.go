package state

import (
	"context"
	"fmt"
	listener "mals/internal/listener"
	client "mals/internal/lsp/client"
	"net"
)

// refactor functions like this
func (s *StateImpl) ListenerListen(listener listener.Listener) {
	s.EventChan <- &EventListenerListen{listener: listener}
}

func (s *StateImpl) ListenerAddConn(listener listener.Listener, ctx context.Context, conn net.Conn) error {
	switch state := state.(type) {
	case *StateListenerLsp:
		client := client.New(s, conn)

		ctx, cancel := context.WithCancel(ctx)

		clientState := NewClientState(cancel)
		state.Clients.Store(client, clientState)

		go func() {
			s.EventChan <- &EventClientLspListen{client: client, ctx: ctx}
		}()

	default:
		s.Warn("ListenerAddConn unhandled listener %v of state type %T", listener, state)
	}

	return nil
}
