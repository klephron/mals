package state

import (
	"context"
	"mals/internal/lsp/client"

	"github.com/puzpuzpuz/xsync/v4"
)

type StateListener interface {
	Cancel()
}

type StateListenerGeneric struct {
	StateListener
	cancel context.CancelFunc
}

func (s *StateListenerGeneric) Cancel() {
	s.cancel()
}

func NewStateListener(cancel context.CancelFunc) *StateListenerGeneric {
	return &StateListenerGeneric{
		cancel: cancel,
	}
}

type StateListenerLsp struct {
	StateListener
	Clients *xsync.Map[*client.Client, *ClientState]
}

func NewStateListenerLsp(cancel context.CancelFunc) *StateListenerLsp {
	return &StateListenerLsp{
		StateListener: NewStateListener(cancel),
		Clients:       xsync.NewMap[*client.Client, *ClientState](),
	}
}

func (s *StateListenerLsp) Cancel() {
	s.Clients.Range(func(key *client.Client, value *ClientState) bool {
		value.Cancel()
		return true
	})
	s.StateListener.Cancel()
}
