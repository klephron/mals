package plane

import (
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"mals/internal/plane/controller/client"
	"mals/internal/plane/controller/lifecycle"
	"mals/internal/plane/controller/listener"
	"mals/internal/plane/controller/logging"
	"mals/internal/plane/event"
	"mals/internal/plane/state"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type Plane struct {
	plane.Plane
	state     *state.State
	bus       *event.EventBus
	client    controller.ClientController
	lifecycle controller.LifecycleController
	listener  controller.ListenerController
	log       controller.LogController
}

func New() plane.Plane {
	plane := &Plane{
		state: &state.State{
			Logs:      xsync.NewMap[string, *state.LogValue](),
			Listeners: xsync.NewMap[string, *state.ListenerValue](),
		},
		bus: event.NewEventBus(),
	}

	plane.client = client.New(plane, plane.state, plane.bus)
	plane.lifecycle = lifecycle.New(plane, plane.state, plane.bus)
	plane.listener = listener.New(plane, plane.state, plane.bus)
	plane.log = logging.New(plane, plane.state, plane.bus)

	return plane
}

func (s *Plane) Client() controller.ClientController {
	return s.client
}

func (s *Plane) Lifecycle() controller.LifecycleController {
	return s.lifecycle
}

func (s *Plane) Listener() controller.ListenerController {
	return s.listener
}

func (s *Plane) Log() controller.LogController {
	return s.log
}

func (s *Plane) Serve(onReady func()) {
	var wg sync.WaitGroup
	var wgReady sync.WaitGroup

	wgReady.Add(4)

	wg.Go(func() {
		s.client.Serve(func() { wgReady.Done() })
	})
	wg.Go(func() {
		s.lifecycle.Serve(func() { wgReady.Done() })
	})
	wg.Go(func() {
		s.listener.Serve(func() { wgReady.Done() })
	})
	wg.Go(func() {
		s.log.Serve(func() { wgReady.Done() })
	})

	wgReady.Wait()
	onReady()

	wg.Wait()
}
