package listener

import (
	"context"
	"fmt"
	"mals/internal/listener"
	"mals/internal/listener/api"
	"mals/internal/listener/lsp"
	"mals/internal/plane"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type state struct {
	statusRW     sync.RWMutex
	statusCancel context.CancelFunc
	listeners    *xsync.Map[string, *stateListener]
}

type stateListener struct {
	rw         sync.RWMutex
	config     *config.Listener
	cancelFunc context.CancelFunc
	mixin      stateListenerM
}

type stateListenerM interface {
	Listener() listener.Listener
}

type stateListenerMApi struct {
	listener *api.ListenerApi
}

func (s *stateListenerMApi) Listener() listener.Listener {
	return s.listener
}

type stateListenerMLsp struct {
	listener *lsp.ListenerLsp
	clients  *xsync.Map[string, *stateListenerLspClient]
}

func (s *stateListenerMLsp) Listener() listener.Listener {
	return s.listener
}

type stateListenerLspClient struct {
	client     listener.ListenerLspClient
	cancelFunc context.CancelFunc
}

type ListenerController struct {
	state state
	plane plane.Plane
}

func New(plane plane.Plane) *ListenerController {
	return &ListenerController{
		state: state{
			statusRW:     sync.RWMutex{},
			statusCancel: nil,
			listeners:    xsync.NewMap[string, *stateListener](),
		},
		plane: plane,
	}
}

func (s *ListenerController) ControllerRun(onReady func()) error {
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

func (s *ListenerController) ControllerShutdown() error {
	s.state.listeners.Range(func(key string, value *stateListener) bool {
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
