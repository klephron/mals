package scope

import (
	"context"
	"fmt"
	"mals/internal/plane"
	"mals/internal/scope"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type ScopeController struct {
	state State
	plane plane.Plane
}

func New(plane plane.Plane) *ScopeController {
	return &ScopeController{
		state: State{
			statusRW:     sync.RWMutex{},
			statusCancel: nil,

			lsps:   xsync.NewMap[string, *config.Lsp](),
			models: xsync.NewMap[string, *config.Model](),

			root: newSpaceRoot(),
		},
		plane: plane,
	}
}

func (s *ScopeController) ControllerRun(onReady func()) error {
	s.state.statusRW.Lock()

	if s.state.statusCancel != nil {
		s.state.statusRW.Unlock()

		err := fmt.Errorf("%T is already serving", s)
		s.plane.Errorf("%v", err)
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

func (s *ScopeController) ControllerShutdown() error {
	s.Close(scope.NewScopeGlobal())

	s.state.statusRW.RLock()
	cancel := s.state.statusCancel
	s.state.statusRW.RUnlock()

	cancel()

	return nil
}
