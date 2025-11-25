package api

import (
	"context"
	"mals/internal/listener/common"
	"mals/internal/state"
)

type ListenerApi struct {
	common.Listener
	state     *state.State
	port      int
	listening bool
}

func Type() string {
	return "lsp"
}

func New(state *state.State, port int) (*ListenerApi, error) {
	l := &ListenerApi{
		state:     state,
		port:      port,
		listening: false,
	}
	return l, nil
}

func (l *ListenerApi) Type() string {
	return Type()
}

func (l *ListenerApi) ListenAndServe(ctx context.Context) error {
	l.listening = true
	err := l.listen(ctx)
	l.listening = false
	return err
}

func (l *ListenerApi) listen(ctx context.Context) error {
	return nil
}

func (l *ListenerApi) Listening() bool {
	return l.listening
}
