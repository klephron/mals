package listener

import (
	"context"
	"fmt"
	"mals/internal/plane"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type ListenerController struct {
	state State
	plane plane.Plane
}

func New(plane plane.Plane) *ListenerController {
	return &ListenerController{
		state: State{
			statusRW:     sync.RWMutex{},
			statusCancel: nil,
			listeners:    xsync.NewMap[string, *Listener](),
		},
		plane: plane,
	}
}

func (s *ListenerController) Run(onReady func()) error {
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
