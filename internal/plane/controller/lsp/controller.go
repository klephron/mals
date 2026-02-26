package lsp

import (
	"context"
	"fmt"
	"mals/internal/plane"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type LspController struct {
	state State
	plane plane.Plane
}

func New(plane plane.Plane) *LspController {
	return &LspController{
		state: State{
			statusRW:     sync.RWMutex{},
			statusCancel: nil,
			lsps:         xsync.NewMap[string, *Lsp](),
		},
		plane: plane,
	}
}

func (s *LspController) ControllerRun(onReady func()) error {
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

func (s *LspController) ControllerShutdown() error {
	s.state.lsps.Range(func(key string, value *Lsp) bool {
		s.Stop(key)
		s.Delete(key)
		return true
	})

	s.state.statusRW.RLock()
	cancel := s.state.statusCancel
	s.state.statusRW.RUnlock()

	cancel()

	return nil
}
