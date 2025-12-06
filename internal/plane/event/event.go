package event

type Event interface {
	event()
}

type EventGeneric struct {
	Event
}

func (*EventGeneric) event() {}
