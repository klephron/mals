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
}

func New(state *state.State, bus *event.EventBus) *LifecycleController {
	return &LifecycleController{
		state:    state,
		bus:      bus,
		external: nil,
	}
}

func (s *LifecycleController) Serve(onReady func()) error {
	if s.external != nil {
		err := fmt.Errorf("%T is already serving", s)

		s.bus.Broadcast(event.EventLog{
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

	onReady()

	for e := range s.external {
		switch e := e.(type) {
		case *event.EventShutdown:
			return nil
		case *event.EventTerminate:
			return nil
		default:
			s.bus.Broadcast(&event.EventLog{
				Level: log.LevelDebug,
				Msg:   fmt.Sprintf("%T unhandled message %T, %v", s, e, e),
			}, s.external)
		}
	}

	return nil
}
