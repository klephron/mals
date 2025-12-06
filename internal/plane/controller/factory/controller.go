package factory

import (
	"errors"
	"mals/internal/plane/event"
	"mals/internal/plane/state"
)

type FactoryController struct {
	state    *state.State
	bus      *event.EventBus
	external <-chan event.Event
}

func NewController(state *state.State, bus *event.EventBus) *FactoryController {
	return &FactoryController{
		state:    state,
		bus:      bus,
		external: nil,
	}
}

func (s *FactoryController) Serve() error {
	if s.external != nil {
		return errors.New("FactoryController is already serving")
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
