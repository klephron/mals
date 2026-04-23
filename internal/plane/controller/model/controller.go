package model

import (
	"context"
	"fmt"
	"mals/internal/model/queued"
	"mals/internal/plane"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type state struct {
	statusRW     sync.RWMutex
	statusCancel context.CancelFunc

	models *xsync.Map[string, *stateModel]
}

type stateModel struct {
	rw         sync.RWMutex
	config     *config.Model
	model      *queued.ModelQueued
	cancelFunc context.CancelFunc
}

type ModelController struct {
	state state
	plane plane.Plane
}

func New(plane plane.Plane) *ModelController {
	return &ModelController{
		state: state{
			statusRW:     sync.RWMutex{},
			statusCancel: nil,
			models:       xsync.NewMap[string, *stateModel](),
		},
		plane: plane,
	}
}

func (s *ModelController) ControllerRun(onReady func()) error {
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

func (s *ModelController) ControllerShutdown() error {
	s.state.models.Range(func(key string, value *stateModel) bool {
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
