package model

import (
	"fmt"
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"mals/internal/plane/event"

	"github.com/puzpuzpuz/xsync/v4"
)

type ModelController struct {
	controller.ModelController

	state State

	plane plane.Plane
	bus   *event.EventBus
}

func New(plane plane.Plane, bus *event.EventBus) *ModelController {
	return &ModelController{
		state: State{
			models:    xsync.NewMap[string, *ModelValue](),
			eventChan: nil,
			taskChan:  make(chan Task),
		},
		plane: plane,
		bus:   bus,
	}
}

func (s *ModelController) Serve(onReady func()) error {
	if s.state.eventChan != nil {
		err := fmt.Errorf("%T is already serving", s)
		s.plane.Log().Errorf("%v", err)
		return err
	}

	s.state.eventChan = s.bus.Subscribe()
	defer func() {
		s.bus.Unsubscribe(s.state.eventChan)
		s.state.eventChan = nil
	}()

	onReady()

	for {
		select {
		case e := <-s.state.eventChan:
			switch e := e.(type) {
			default:
				s.plane.Log().Warnf("%T unhandled message %T, %v", s, e, e)
			}

		case t := <-s.state.taskChan:
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
