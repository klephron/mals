package listener

import (
	"context"
	"mals/internal/log"
)

type ListenerApi struct {
	Listener
	port int
	log  *log.Log
}

func listenerApiType() string {
	return "lsp"
}

func newListenerApi(port int, log *log.Log) (*ListenerApi, error) {
	l := &ListenerApi{
		port: port,
		log:  log,
	}
	return l, nil
}

func (l *ListenerApi) Type() string {
	return listenerApiType()
}

func (*ListenerApi) ListenAndServe(ctx context.Context) error {
	return nil
}
