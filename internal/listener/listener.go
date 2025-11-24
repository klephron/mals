package listener

import (
	"context"
	"fmt"
	"mals/internal/log"
	"mals/pkg/config"
)

type Listener interface {
	Type() string
	ListenAndServe(ctx context.Context) error
}

func NewListener(listener *config.Listener, log *log.Log) (Listener, error) {
	switch listener.Type {
	case listenerLspType():
		return newListenerLsp(listener.Port, log)
	case listenerApiType():
		return newListenerApi(listener.Port, log)
	default:
		return nil, fmt.Errorf(`unhandled listener type "%v"`, listener.Type)
	}
}
