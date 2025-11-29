package state

import (
	"context"
	"mals/internal/lsp/client"

	"github.com/puzpuzpuz/xsync/v4"
)

type ListenerValue interface {
	Cancel()
}

type ListenerValueLsp struct {
	ListenerValue
	CancelFunc context.CancelFunc
	Clients    *xsync.Map[*client.Client, struct{}]
}

func NewListenerValueLsp(cancel context.CancelFunc) *ListenerValueLsp {
	return &ListenerValueLsp{
		CancelFunc: cancel,
		Clients:    xsync.NewMap[*client.Client, struct{}](),
	}
}
