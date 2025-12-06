package event

import listener "mals/internal/listener"

type EventListenerAdd struct {
	EventGeneric
	Listener listener.Listener
}

type EventListenerDelete struct {
	EventGeneric
	Listener listener.Listener
}

type EventListenerStart struct {
	EventGeneric
	Listener listener.Listener
}

type EventListenerStop struct {
	EventGeneric
	Listener listener.Listener
}
