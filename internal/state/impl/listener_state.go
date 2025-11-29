package state

import (
	"context"
	"mals/internal/lsp/client"

	"github.com/puzpuzpuz/xsync/v4"
)

type ListenerState interface {
	Cancel()
}

type ListenerStateGeneric struct {
	ListenerState
	cancel context.CancelFunc
}

func (s *ListenerStateGeneric) Cancel() {
	s.cancel()
}

func NewListenerStateGeneric(cancel context.CancelFunc) *ListenerStateGeneric {
	return &ListenerStateGeneric{
		cancel: cancel,
	}
}

type ListenerStateLsp struct {
	ListenerState
	Clients *xsync.Map[*client.Client, *ClientState]
}

func NewListenerStateLsp(cancel context.CancelFunc) *ListenerStateLsp {
	return &ListenerStateLsp{
		ListenerState: NewListenerStateGeneric(cancel),
		Clients:       xsync.NewMap[*client.Client, *ClientState](),
	}
}

func (s *ListenerStateLsp) Cancel() {
	s.Clients.Range(func(key *client.Client, value *ClientState) bool {
		key.Close()
		return true
	})
	s.ListenerState.Cancel()
}
