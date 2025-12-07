package logging

import (
	"mals/internal/log"
	"mals/pkg/config"
)

type Event interface {
	event()
}

type EventGeneric struct {
	Event
	Result chan error
}

func (*EventGeneric) event() {}

func NewEventSingle() EventGeneric {
	return EventGeneric{Result: make(chan error, 1)}
}

type EventRegister struct {
	EventGeneric
	Config config.Log
}

type EventUnregister struct {
	EventGeneric
	Name string
}

type EventCreate struct {
	EventGeneric
	Name string
}

type EventDelete struct {
	EventGeneric
	Name string
}

type EventStart struct {
	EventGeneric
	Name string
}

type EventStop struct {
	EventGeneric
	Name string
}

type EventLog struct {
	EventGeneric
	Level log.Level
	Msg   string
}
