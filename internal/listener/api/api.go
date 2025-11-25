package api

import (
	"context"
	"mals/internal/listener/common"
	"mals/internal/state"
)

type ListenerApi struct {
	common.Listener
	state *state.State
	port  int
}

func Type() string {
	return "lsp"
}

func New(state *state.State, port int) (*ListenerApi, error) {
	l := &ListenerApi{
		state: state,
		port:  port,
	}
	return l, nil
}

func (l *ListenerApi) Type() string {
	return Type()
}

func (*ListenerApi) ListenAndServe(ctx context.Context) error {
	return nil
}
