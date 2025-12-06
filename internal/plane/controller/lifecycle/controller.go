package lifecycle

import (
	"fmt"
	"mals/internal/plane/event"
	"mals/internal/plane/state"
)

type LifecycleController struct {
	state    *state.State
	bus      *event.EventBus
	external <-chan event.Event
}

func NewController(state *state.State, bus *event.EventBus) *LifecycleController {
	return &LifecycleController{
		state:    state,
		bus:      bus,
		external: nil,
	}
}

func (s *LifecycleController) Serve() error {
	if s.external != nil {
		return fmt.Errorf("%T is already serving", s)
	}

	s.external = s.bus.Subscribe()
	defer func() {
		s.bus.Unsubscribe(s.external)
		s.external = nil
	}()

	for e := range s.external {
		switch e := e.(type) {
		case *event.EventShutdown:
			s.handleShutdown(e)
		case *event.EventTerminate:
			s.handleTerminate(e)
			return nil
		}
	}

	return nil
}
