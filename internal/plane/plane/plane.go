package plane

import (
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"mals/internal/plane/controller/client"
	"mals/internal/plane/controller/listener"
	"mals/internal/plane/controller/logging"
	"mals/internal/plane/controller/model"
	"mals/internal/plane/controller/scope"
	"mals/internal/plane/controller/usage"
	"sync"
)

type Plane struct {
	plane.Plane
	client   controller.ClientController
	listener controller.ListenerController
	log      controller.LogController
	model    controller.ModelController
	scope    controller.ScopeController
	usage    controller.UsageController
}

func New() plane.Plane {
	plane := &Plane{}

	plane.client = client.New(plane)
	plane.listener = listener.New(plane)
	plane.log = logging.New(plane)
	plane.model = model.New(plane)
	plane.scope = scope.New(plane)
	plane.usage = usage.New(plane)

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

func (s *Plane) Scope() controller.ScopeController {
	return s.scope
}

func (s *Plane) Usage() controller.UsageController {
	return s.usage
}

func (s *Plane) Serve(onReady func()) {
	var wg sync.WaitGroup
	var wgReady sync.WaitGroup

	{
		wgReady.Add(1)
		wg.Go(func() {
			s.log.Serve(func() { wgReady.Done() })
		})
	}
	{
		wgReady.Add(1)
		wg.Go(func() {
			s.model.Serve(func() { wgReady.Done() })
		})
	}
	{
		wgReady.Add(1)
		wg.Go(func() {
			s.scope.Serve(func() { wgReady.Done() })
		})
	}
	{
		wgReady.Add(1)
		wg.Go(func() {
			s.usage.Serve(func() { wgReady.Done() })
		})
	}

	{
		wgReady.Add(1)
		wg.Go(func() {
			s.client.Serve(func() { wgReady.Done() })
		})
	}
	{
		wgReady.Add(1)
		wg.Go(func() {
			s.listener.Serve(func() { wgReady.Done() })
		})
	}

	wgReady.Wait()
	onReady()
	wg.Wait()
}

func (s *Plane) Shutdown() error {
	if err := s.client.Shutdown(); err != nil {
		s.Log().Errorf("%v", err)
		return err
	}
	if err := s.listener.Shutdown(); err != nil {
		s.Log().Errorf("%v", err)
		return err
	}

	if err := s.usage.Shutdown(); err != nil {
		s.Log().Errorf("%v", err)
		return err
	}
	if err := s.scope.Shutdown(); err != nil {
		s.Log().Errorf("%v", err)
		return err
	}
	if err := s.model.Shutdown(); err != nil {
		s.Log().Errorf("%v", err)
		return err
	}
	if err := s.log.Shutdown(); err != nil {
		s.Log().Errorf("%v", err)
		return err
	}
	return nil
}
