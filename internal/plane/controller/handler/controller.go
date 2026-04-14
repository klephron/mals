package handler

import (
	"context"
	"fmt"
	"mals/internal/plane"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type HandlerController struct {
	state State
	plane plane.Plane
}

func New(plane plane.Plane) *HandlerController {
	return &HandlerController{
		state: State{
			statusRW:     sync.RWMutex{},
			statusCancel: nil,
			handlers:     xsync.NewMap[string, *config.Handler](),
		},
		plane: plane,
	}
}

func (s *HandlerController) ControllerRun(onReady func()) error {
	s.state.statusRW.Lock()

	if s.state.statusCancel != nil {
		s.state.statusRW.Unlock()

		err := fmt.Errorf("%T is already serving", s)
		s.plane.Errorf("%T: %v", s, err)
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.state.statusCancel = cancel
	s.state.statusRW.Unlock()

	onReady()
	<-ctx.Done()

	s.state.statusRW.Lock()
	s.state.statusCancel = nil
	s.state.statusRW.Unlock()

	return nil
}

func (s *HandlerController) ControllerShutdown() error {
	s.state.statusRW.RLock()
	cancel := s.state.statusCancel
	s.state.statusRW.RUnlock()

	cancel()

	return nil
}
