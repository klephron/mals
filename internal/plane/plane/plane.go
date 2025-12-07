package plane

import (
	c "mals/internal/client"
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"mals/internal/plane/controller/client"
	"mals/internal/plane/controller/listener"
	"mals/internal/plane/controller/logging"
	"mals/internal/plane/event"
	"mals/internal/plane/state"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type Plane struct {
	plane.Plane
	state    *state.State
	bus      *event.EventBus
	client   controller.ClientController
	listener controller.ListenerController
	log      controller.LogController
}

func New() plane.Plane {
	plane := &Plane{
		state: &state.State{
			Logs:      xsync.NewMap[string, *state.LogValue](),
			Listeners: xsync.NewMap[string, *state.ListenerValue](),
			Clients:   xsync.NewMap[c.Client, *state.ClientValue](),
		},
		bus: event.NewEventBus(),
	}

	plane.client = client.New(plane, plane.state, plane.bus)
	plane.listener = listener.New(plane, plane.state, plane.bus)
	plane.log = logging.New(plane, plane.state, plane.bus)

	return plane
}

func (s *Plane) Client() controller.ClientController {
	return s.client
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

	wgReady.Add(3)

	wg.Go(func() {
		s.client.Serve(func() { wgReady.Done() })
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

func (s *Plane) Shutdown() error {
	if err := s.listener.Shutdown(); err != nil {
		return err
	}
	if err := s.client.Shutdown(); err != nil {
		return err
	}
	if err := s.log.Shutdown(); err != nil {
		return err
	}
	return nil
}

func (s *Plane) Terminate() error {
	s.listener.Terminate()
	s.client.Terminate()
	s.log.Terminate()
	return nil
}
