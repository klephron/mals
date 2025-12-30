package client

import (
	"context"
	"fmt"
	"mals/internal/client"
	"mals/internal/plane"
	"mals/internal/plane/controller"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type ClientController struct {
	controller.ClientController
	state State
	plane plane.Plane
}

func New(plane plane.Plane) *ClientController {
	return &ClientController{
		state: State{
			statusRW:     sync.RWMutex{},
			statusCancel: nil,
			clients:      xsync.NewMap[client.Client, *ClientValue](),
		},
		plane: plane,
	}
}

func (s *ClientController) Run(onReady func()) error {
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
