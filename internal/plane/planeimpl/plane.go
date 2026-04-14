package planeimpl

import (
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"mals/internal/plane/controller/handler"
	"mals/internal/plane/controller/listener"
	"mals/internal/plane/controller/log"
	"mals/internal/plane/controller/lsp"
	"mals/internal/plane/controller/model"
	"mals/internal/plane/controller/scope"
	"sync"
)

type Plane struct {
	listener controller.ListenerController
	log      controller.LogController
	lsp      controller.LspController
	model    controller.ModelController
	scope    controller.ScopeController
	handler  controller.HandlerController
}

func New() plane.Plane {
	plane := &Plane{}

	plane.listener = listener.New(plane)
	plane.log = log.New(plane)
	plane.lsp = lsp.New(plane)
	plane.model = model.New(plane)
	plane.scope = scope.New(plane)
	plane.handler = handler.New(plane)

	return plane
}

func (s *Plane) Run(onReady func()) {
	var wg sync.WaitGroup
	var wgReady sync.WaitGroup

	{
		wgReady.Add(1)
		wg.Go(func() {
			s.log.ControllerRun(func() { wgReady.Done() })
		})
	}
	{
		wgReady.Add(1)
		wg.Go(func() {
			s.lsp.ControllerRun(func() { wgReady.Done() })
		})
	}
	{
		wgReady.Add(1)
		wg.Go(func() {
			s.model.ControllerRun(func() { wgReady.Done() })
		})
	}
	{
		wgReady.Add(1)
		wg.Go(func() {
			s.scope.ControllerRun(func() { wgReady.Done() })
		})
	}
	{
		wgReady.Add(1)
		wg.Go(func() {
			s.handler.ControllerRun(func() { wgReady.Done() })
		})
	}

	{
		wgReady.Add(1)
		wg.Go(func() {
			s.listener.ControllerRun(func() { wgReady.Done() })
		})
	}

	wgReady.Wait()
	onReady()
	wg.Wait()
}

func (s *Plane) Shutdown() error {
	if err := s.listener.ControllerShutdown(); err != nil {
		s.Errorf("%v", err)
		return err
	}

	if err := s.handler.ControllerShutdown(); err != nil {
		s.Errorf("%v", err)
		return err
	}
	if err := s.scope.ControllerShutdown(); err != nil {
		s.Errorf("%v", err)
		return err
	}
	if err := s.model.ControllerShutdown(); err != nil {
		s.Errorf("%v", err)
		return err
	}
	if err := s.lsp.ControllerShutdown(); err != nil {
		s.Errorf("%v", err)
		return err
	}
	if err := s.log.ControllerShutdown(); err != nil {
		s.Errorf("%v", err)
		return err
	}
	return nil
}
