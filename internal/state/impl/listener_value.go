package state

import (
	"context"
	"mals/internal/lsp/client"

	"github.com/puzpuzpuz/xsync/v4"
)

type ListenerValue interface {
	Cancel()
}

type ListenerValueGeneric struct {
	ListenerValue
	cancel context.CancelFunc
}

func (s *ListenerValueGeneric) Cancel() {
	s.cancel()
}

func NewListenerValueGeneric(cancel context.CancelFunc) *ListenerValueGeneric {
	return &ListenerValueGeneric{
		cancel: cancel,
	}
}

type ListenerValueLsp struct {
	ListenerValue
	Clients *xsync.Map[*client.Client, struct{}]
}

func NewListenerValueLsp(cancel context.CancelFunc) *ListenerValueLsp {
	return &ListenerValueLsp{
		ListenerValue: NewListenerValueGeneric(cancel),
		Clients:       xsync.NewMap[*client.Client, struct{}](),
	}
}
