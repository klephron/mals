package controller

import (
	"mals/internal/client"
	"mals/pkg/config"
)

type UsageController interface {
	Shutdown() error
	Serve(onReady func()) error

	RegisterModel(config config.Model) error
	RegisterLsp(config config.Lsp) error
	RegisterUsage(config config.Usage) error

	ClientSubscribe(client client.Client)
	ClientUnsubscribe(client client.Client)
}
