package listener

import (
	"fmt"
	"mals/internal/listener/api"
	"mals/internal/listener/common"
	"mals/internal/listener/lsp"
	"mals/internal/state"
	"mals/pkg/config"
)

func New(state *state.State, listener *config.Listener) (common.Listener, error) {
	switch listener.Type {
	case lsp.Type():
		return lsp.New(state, listener.Port)
	case api.Type():
		return lsp.New(state, listener.Port)
	default:
		return nil, fmt.Errorf(`unhandled listener type "%v"`, listener.Type)
	}
}
