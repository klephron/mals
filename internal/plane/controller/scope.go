package controller

import (
	"mals/pkg/config"
)

type ScopeController interface {
	Shutdown() error
	Serve(onReady func()) error

	RegisterModel(config config.Model) error
	RegisterLsp(config config.Lsp) error
}
