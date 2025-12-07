package listener

import "mals/pkg/config"

type Event interface {
	event()
}

type EventGeneric struct {
	Event
}

func (*EventGeneric) event() {}

type EventRegister struct {
	EventGeneric
	Config config.Listener

	Error error
}

type EventUnregister struct {
	EventGeneric
	Name string

	Error error
}

type EventCreate struct {
	EventGeneric
	Name string

	Error error
}

type EventDelete struct {
	EventGeneric
	Name string

	Error error
}

type EventStart struct {
	EventGeneric
	Name string

	Error error
}

type EventStop struct {
	EventGeneric
	Name string

	Error error
}
