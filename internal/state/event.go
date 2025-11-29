package state

type Event interface {
	event()
}

type EventListenerDone struct {
	Event
}

func (*EventListenerDone) event() {
}
