package lsp

import (
	"context"
	"fmt"
	"mals/internal/lsp/server"
	"mals/internal/plane"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type state struct {
	statusRW     sync.RWMutex
	statusCancel context.CancelFunc

	lsps *xsync.Map[string, *stateLsp]
}

type stateLsp struct {
	rw         sync.RWMutex
	config     *config.Lsp
	lsp        server.LspServer
	cancelFunc context.CancelFunc
}

type LspController struct {
	state state
	plane plane.Plane
}

func New(plane plane.Plane) *LspController {
	return &LspController{
		state: state{
			statusRW:     sync.RWMutex{},
			statusCancel: nil,
			lsps:         xsync.NewMap[string, *stateLsp](),
		},
		plane: plane,
	}
}

func (s *LspController) ControllerRun(onReady func()) error {
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

func (s *LspController) ControllerShutdown() error {
	s.state.lsps.Range(func(key string, value *stateLsp) bool {
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
