package listener

import (
	"fmt"
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"mals/internal/plane/event"

	"github.com/puzpuzpuz/xsync/v4"
)

type ListenerController struct {
	controller.ListenerController

	state State

	plane plane.Plane
	bus   *event.EventBus
}

func New(plane plane.Plane, bus *event.EventBus) *ListenerController {
	return &ListenerController{
		state: State{
			listeners: xsync.NewMap[string, *ListenerValue](),
			external:  nil,
			internal:  make(chan Task),
		},
		plane: plane,
		bus:   bus,
	}
}

func (s *ListenerController) Serve(onReady func()) error {
	if s.state.external != nil {
		err := fmt.Errorf("%T is already serving", s)
		s.plane.Log().Errorf("%v", err)
		return err
	}

	s.state.external = s.bus.Subscribe()
	defer func() {
		s.bus.Unsubscribe(s.state.external)
		s.state.external = nil
	}()

	onReady()

	for {
		select {
		case e := <-s.state.external:
			switch e := e.(type) {
			default:
				s.plane.Log().Warnf("%T unhandled message %T, %v", s, e, e)
			}

		case t := <-s.state.internal:
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
			case *TaskClientAdd:
				s.handleClientAdd(t)
			case *TaskClientRemove:
				s.handleClientRemove(t)
			default:
				s.plane.Log().Warnf("%T unhandled internal message %T, %v", s, t, t)
			}
		}
	}
}
