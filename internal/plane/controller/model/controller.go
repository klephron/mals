package model

import (
	"fmt"
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"mals/internal/plane/event"
	"mals/internal/plane/state"
)

type ModelController struct {
	controller.ModelController
	plane    plane.Plane
	state    *state.State
	bus      *event.EventBus
	external <-chan event.Event
	internal chan Task
}

func New(plane plane.Plane, state *state.State, bus *event.EventBus) *ModelController {
	return &ModelController{
		plane:    plane,
		state:    state,
		bus:      bus,
		external: nil,
		internal: make(chan Task),
	}
}

func (s *ModelController) Serve(onReady func()) error {
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
			case *TaskShutdown:
				s.handleShutdown(t)
				return nil
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
