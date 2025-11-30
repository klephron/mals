package api

import (
	"context"
	"mals/internal/control/controller"
	"mals/internal/listener"
)

type ListenerApi struct {
	listener.Listener
	controller *controller.Controller
	port       int
	listening  bool
}

func Type() string {
	return "api"
}

func New(controller *controller.Controller, port int) (*ListenerApi, error) {
	l := &ListenerApi{
		controller: controller,
		port:       port,
		listening:  false,
	}
	return l, nil
}

func (s *ListenerApi) Type() string {
	return Type()
}

func (s *ListenerApi) Listen(ctx context.Context) error {
	s.listening = true
	err := s.listen(ctx)
	s.listening = false
	return err
}

func (s *ListenerApi) listen(ctx context.Context) error {
	return nil
}

func (s *ListenerApi) Listening() bool {
	return s.listening
}
