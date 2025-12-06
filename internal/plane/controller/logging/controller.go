package logging

import (
	"fmt"
	"mals/internal/plane/controller"
	"mals/internal/plane/event"
	"mals/internal/plane/state"
)

type LogController struct {
	controller.LogController
	state    *state.State
	bus      *event.EventBus
	external <-chan event.Event
}

func New(state *state.State, bus *event.EventBus) *LogController {
	return &LogController{
		state:    state,
		bus:      bus,
		external: nil,
	}
}

func (s *LogController) Serve() error {
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
			return nil
		case *event.EventTerminate:
			return nil
		}
	}

	return nil
}
