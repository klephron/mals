package listener

import (
	"fmt"
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"mals/internal/plane/event"
	"mals/internal/plane/state"
)

type ListenerController struct {
	controller.ListenerController
	plane    plane.Plane
	state    *state.State
	bus      *event.EventBus
	external <-chan event.Event
	internal chan Event
}

func New(plane plane.Plane, state *state.State, bus *event.EventBus) *ListenerController {
	return &ListenerController{
		plane:    plane,
		state:    state,
		bus:      bus,
		external: nil,
		internal: make(chan Event),
	}
}

func (s *ListenerController) Serve(onReady func()) error {
	if s.external != nil {
		err := fmt.Errorf("%T is already serving", s)
		s.plane.Log().Errorf("%v", err)
		return err
	}

	s.external = s.bus.Subscribe()
	defer func() {
		s.bus.Unsubscribe(s.external)
		s.external = nil
	}()

	onReady()

	for {
		select {
		case e := <-s.external:
			switch e := e.(type) {
			case *event.EventShutdown:
				s.handleShutdown(e)
				return nil
			case *event.EventTerminate:
				s.handleTerminate(e)
				return nil
			default:
				s.plane.Log().Warnf("%T unhandled message %T, %v", s, e, e)
			}

		case e := <-s.internal:
			switch e := e.(type) {
			case *EventRegister:
				s.handleRegister(e)
			case *EventUnregister:
				s.handleUnregister(e)
			case *EventCreate:
				s.handleCreate(e)
			case *EventDelete:
				s.handleDelete(e)
			case *EventStart:
				s.handleStart(e)
			case *EventStop:
				s.handleStop(e)
			default:
				s.plane.Log().Warnf("%T unhandled internal message %T, %v", s, e, e)
			}
		}
	}
}
