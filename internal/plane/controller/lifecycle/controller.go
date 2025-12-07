package lifecycle

import (
	"fmt"
	"mals/internal/plane"
	"mals/internal/plane/event"
	"mals/internal/plane/state"
)

type LifecycleController struct {
	plane    plane.Plane
	state    *state.State
	bus      *event.EventBus
	external <-chan event.Event
}

func New(plane plane.Plane, state *state.State, bus *event.EventBus) *LifecycleController {
	return &LifecycleController{
		plane:    plane,
		state:    state,
		bus:      bus,
		external: nil,
	}
}

func (s *LifecycleController) Serve(onReady func()) error {
	if s.external != nil {
		err := fmt.Errorf("%T is already serving", s)
		s.plane.Log().Errorf("%v", err)
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
			s.plane.Log().Warnf("%T unhandled message %T, %v", s, e, e)
		}
	}

	return nil
}
