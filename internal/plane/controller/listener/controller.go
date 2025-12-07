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
	internal chan Task
}

func New(plane plane.Plane, state *state.State, bus *event.EventBus) *ListenerController {
	return &ListenerController{
		plane:    plane,
		state:    state,
		bus:      bus,
		external: nil,
		internal: make(chan Task),
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
			default:
				s.plane.Log().Warnf("%T unhandled message %T, %v", s, e, e)
			}

		case t := <-s.internal:
			switch t := t.(type) {
			case *TaskRegister:
				s.handleRegister(t)
			case *TaskUnregister:
				s.handleUnregister(t)
			case *TaskCreate:
				s.handleCreate(t)
			case *TaskDelete:
				s.handleDelete(t)
			case *TaskStart:
				s.handleStart(t)
			case *TaskStop:
				s.handleStop(t)
			default:
				s.plane.Log().Warnf("%T unhandled internal message %T, %v", s, t, t)
			}
		}
	}
}
