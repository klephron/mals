package lifecycle

type Event interface {
	event()
}

type EventGeneric struct {
	Event
}

func (*EventGeneric) event() {}

type EventShutdown struct {
	EventGeneric
}

type EventTerminate struct {
	EventGeneric
}
