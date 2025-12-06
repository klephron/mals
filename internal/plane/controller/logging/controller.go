package logging

import (
	"fmt"
	"mals/internal/log"
	"mals/internal/plane/controller"
	"mals/internal/plane/event"
	"mals/internal/plane/state"
)

type LogController struct {
	controller.LogController
	state    *state.State
	bus      *event.EventBus
	external <-chan event.Event
	internal chan Event
}

func New(state *state.State, bus *event.EventBus) *LogController {
	return &LogController{
		state:    state,
		bus:      bus,
		external: nil,
		internal: nil,
	}
}

func (s *LogController) Serve() error {
	if s.external != nil {
		err := fmt.Errorf("%T is already serving", s)

		s.bus.Broadcast(event.EventLog{
			Level:   log.LevelError,
			Pattern: "%v",
			Args:    []any{err},
		}, s.external)

		return err
	}

	s.external = s.bus.Subscribe()
	defer func() {
		s.bus.Unsubscribe(s.external)
		s.external = nil
	}()

	s.internal = make(chan Event)
	defer func() {
		close(s.internal)
		s.internal = nil
	}()

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
				s.bus.Broadcast(&event.EventLog{
					Level:   log.LevelWarn,
					Pattern: "%T unhandled message %T, %v",
					Args:    []any{s, e, e},
				}, s.external)
			}

		case e := <-s.internal:
			switch e := e.(type) {
			case *EventRegister:
				s.handleRegister(e)
			case *EventCreate:
				s.handleCreate(e)
			case *EventDelete:
				s.handleDelete(e)
			case *EventStart:
				s.handleStart(e)
			case *EventStop:
				s.handleStop(e)
			default:
				s.bus.Broadcast(&event.EventLog{
					Level:   log.LevelWarn,
					Pattern: "%T unhandled internal message %T, %v",
					Args:    []any{s, e, e},
				}, s.external)
			}
		}
	}
}
