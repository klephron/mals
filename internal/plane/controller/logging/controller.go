package logging

import (
	"context"
	"fmt"
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type LogController struct {
	controller.LogController
	state State
	plane plane.Plane
}

func New(plane plane.Plane) *LogController {
	return &LogController{
		state: State{
			statusRW:     sync.RWMutex{},
			statusCancel: nil,
			logs:         xsync.NewMap[string, *LogValue](),
		},
		plane: plane,
	}
}

func (s *LogController) Run(onReady func()) error {
	s.state.statusRW.Lock()

	if s.state.statusCancel != nil {
		s.state.statusRW.Unlock()

		err := fmt.Errorf("%T is already serving", s)
		s.plane.Log().Errorf("%v", err)
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
