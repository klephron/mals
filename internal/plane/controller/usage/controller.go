package usage

import (
	"context"
	"fmt"
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type UsageController struct {
	controller.UsageController
	state State
	plane plane.Plane
}

func New(plane plane.Plane) *UsageController {
	return &UsageController{
		state: State{
			statusRW:     sync.RWMutex{},
			statusCancel: nil,
			usages:       xsync.NewMap[string, *config.Usage](),
		},
		plane: plane,
	}
}

func (s *UsageController) Run(onReady func()) error {
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
