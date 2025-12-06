package lifecycle

import (
	"fmt"
	"mals/internal/log"
	"mals/internal/plane/event"
	"mals/internal/plane/state"
)

type LifecycleController struct {
	state    *state.State
	bus      *event.EventBus
	external <-chan event.Event
	internal chan Event
}

func New(state *state.State, bus *event.EventBus) *LifecycleController {
	return &LifecycleController{
		state:    state,
		bus:      bus,
		external: nil,
		internal: nil,
	}
}

func (s *LifecycleController) Serve() error {
	if s.external != nil || s.internal != nil {
		return fmt.Errorf("%T is already serving", s)
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
				return nil
			case *event.EventTerminate:
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
			case *EventShutdown:
				s.handleLifecycleShutdown(e)
			case *EventTerminate:
				s.handleLifecycleTerminate(e)
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
