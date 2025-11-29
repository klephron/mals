package state

import (
	"context"
	client "mals/internal/lsp/client"
)

type Event interface {
	event()
}

type EventGeneric struct {
	Event
}

func (*EventGeneric) event() {}

type EventListenerDown struct {
	EventGeneric
}

type EventClientListen struct {
	EventGeneric
	client *client.Client
	ctx    context.Context
}
