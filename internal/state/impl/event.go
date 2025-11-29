package state

import (
	"context"
	listener "mals/internal/listener"
	client "mals/internal/lsp/client"
)

type Event interface {
	event()
}

type EventGeneric struct {
	Event
}

func (*EventGeneric) event() {}

type EventListenerListen struct {
	EventGeneric
	listener listener.Listener
	ctx      context.Context
}

type EventClientLspListen struct {
	EventGeneric
	client *client.Client
	ctx    context.Context
}

type EventShutdown struct {
	EventGeneric
}
