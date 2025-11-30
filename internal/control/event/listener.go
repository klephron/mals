package event

import listener "mals/internal/listener"

type EventListenerListen struct {
	EventGeneric
	listener listener.Listener
}
