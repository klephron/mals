package model

import (
	"context"
	"fmt"
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type ModelController struct {
	controller.ModelController
	state State
	plane plane.Plane
}

func New(plane plane.Plane) *ModelController {
	return &ModelController{
		state: State{
			statusRW:     sync.RWMutex{},
			statusCancel: nil,
			models:       xsync.NewMap[string, *ModelValue](),
		},
		plane: plane,
	}
}

func (s *ModelController) Run(onReady func()) error {
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

	// s.plane.Log().Infof("%T done", s)

	s.state.statusRW.Lock()
	s.state.statusCancel = nil
	s.state.statusRW.Unlock()

	return nil
}
