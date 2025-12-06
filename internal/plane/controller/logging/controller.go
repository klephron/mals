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

func (s *LogController) Serve(onReady func()) error {
	if s.external != nil {
		err := fmt.Errorf("%T is already serving", s)

		s.bus.Unicast(event.EventLog{
			Level: log.LevelError,
			Msg:   fmt.Sprintf("%v", err),
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
			case *event.EventLog:
				s.handleLog(e)
			default:
				s.bus.Unicast(&event.EventLog{
					Level: log.LevelDebug,
					Msg:   fmt.Sprintf("%T unhandled message %T, %v", s, e, e),
				}, s.external)
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
				s.bus.Unicast(&event.EventLog{
					Level: log.LevelWarn,
					Msg:   fmt.Sprintf("%T unhandled internal message %T, %v", s, e, e),
				}, s.external)
			}
		}
	}
}
