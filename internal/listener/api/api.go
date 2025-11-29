package api

import (
	"context"
	"mals/internal/listener"
	"mals/internal/state"
)

type ListenerApi struct {
	listener.Listener
	state     state.State
	port      int
	listening bool
}

func Type() string {
	return "lsp"
}

func New(state state.State, port int) (*ListenerApi, error) {
	l := &ListenerApi{
		state:     state,
		port:      port,
		listening: false,
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
