package scope

import (
	"context"
	"fmt"
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type ScopeController struct {
	controller.ScopeController
	state State
	plane plane.Plane
}

func New(plane plane.Plane) *ScopeController {
	return &ScopeController{
		state: State{
			statusRW:     sync.RWMutex{},
			statusCancel: nil,

			models: xsync.NewMap[string, *config.Model](),

			root: newSpaceRoot(),
		},
		plane: plane,
	}
}

func (s *ScopeController) Run(onReady func()) error {
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

	// s.plane.Log().Infof("%T done", s)

	s.state.statusRW.Lock()
	s.state.statusCancel = nil
	s.state.statusRW.Unlock()

	return nil
}
