package log

import (
	"context"
	"fmt"
	"mals/internal/log"
	"mals/internal/plane"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type state struct {
	statusCancel context.CancelFunc
	statusRW     sync.RWMutex

	logs *xsync.Map[string, *stateLog]
}

type stateLog struct {
	rw      sync.RWMutex
	config  *config.Log
	log     log.Log
	enabled bool
}

type LogController struct {
	state state
	plane plane.Plane
}

func New(plane plane.Plane) *LogController {
	return &LogController{
		state: state{
			statusRW:     sync.RWMutex{},
			statusCancel: nil,
			logs:         xsync.NewMap[string, *stateLog](),
		},
		plane: plane,
	}
}

func (s *LogController) ControllerRun(onReady func()) error {
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

func (s *LogController) ControllerShutdown() error {
	s.state.logs.Range(func(key string, value *stateLog) bool {
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
