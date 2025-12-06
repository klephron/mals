package ownership

import (
	"errors"
	"mals/internal/plane/event"
	"mals/internal/plane/state"
)

type OwnershipController struct {
	state    *state.State
	bus      *event.EventBus
	external <-chan event.Event
}

func NewController(state *state.State, bus *event.EventBus) *OwnershipController {
	return &OwnershipController{
		state:    state,
		bus:      bus,
		external: nil,
	}
}

func (s *OwnershipController) Serve() error {
	if s.external != nil {
		return errors.New("OwnershipController is already serving")
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
