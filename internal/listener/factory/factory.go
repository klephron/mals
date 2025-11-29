package listener

import (
	"fmt"
	"mals/internal/listener"
	"mals/internal/listener/api"
	"mals/internal/listener/lsp"
	"mals/internal/state"
	"mals/pkg/config"
)

func NewConfig(state state.State, listener *config.Listener) (listener.Listener, error) {
	switch listener.Type {
	case lsp.Type():
		return lsp.New(state, listener.Port)
	case api.Type():
		return api.New(state, listener.Port)
	default:
		return nil, fmt.Errorf(`unhandled listener type "%v"`, listener.Type)
	}
}
