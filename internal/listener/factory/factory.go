package listener

import (
	"fmt"
	"mals/internal/control/controller"
	"mals/internal/listener"
	"mals/internal/listener/api"
	"mals/internal/listener/lsp"
	"mals/pkg/config"
)

func NewConfig(controller *controller.Controller, listener *config.Listener) (listener.Listener, error) {
	switch listener.Type {
	case lsp.Type():
		return lsp.New(controller, listener.Port)
	case api.Type():
		return api.New(controller, listener.Port)
	default:
		return nil, fmt.Errorf(`unhandled listener type "%v"`, listener.Type)
	}
}
