package plane

import (
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"mals/internal/plane/controller/client"
	"mals/internal/plane/controller/listener"
	"mals/internal/plane/controller/logging"
	"mals/internal/plane/controller/model"
	"mals/internal/plane/event"
	"sync"
)

type Plane struct {
	plane.Plane
	bus      *event.EventBus
	client   controller.ClientController
	listener controller.ListenerController
	log      controller.LogController
	model    controller.ModelController
}

func New() plane.Plane {
	plane := &Plane{
		bus: event.NewEventBus(),
	}

	plane.client = client.New(plane)
	plane.listener = listener.New(plane, plane.bus)
	plane.log = logging.New(plane, plane.bus)
	plane.model = model.New(plane, plane.bus)

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

func (s *Plane) Model() controller.ModelController {
	return s.model
}

func (s *Plane) Serve(onReady func()) {
	var wg sync.WaitGroup
	var wgReady sync.WaitGroup

	wgReady.Add(1)
	{
		wg.Go(func() {
			s.client.Serve(func() { wgReady.Done() })
		})
		// wg.Go(func() {
		// 	s.listener.Serve(func() { wgReady.Done() })
		// })
		// wg.Go(func() {
		// 	s.log.Serve(func() { wgReady.Done() })
		// })
		// wg.Go(func() {
		// 	s.model.Serve(func() { wgReady.Done() })
		// })
	}
	wgReady.Wait()

	onReady()

	wg.Wait()
}

func (s *Plane) Shutdown() error {
	// if err := s.model.Shutdown(); err != nil {
	// 	s.Log().Errorf("%v", err)
	// 	return err
	// }
	// if err := s.listener.Shutdown(); err != nil {
	// 	s.Log().Errorf("%v", err)
	// 	return err
	// }
	if err := s.client.Shutdown(); err != nil {
		s.Log().Errorf("%v", err)
		return err
	}
	// if err := s.log.Shutdown(); err != nil {
	// 	s.Log().Errorf("%v", err)
	// 	return err
	// }
	return nil
}
