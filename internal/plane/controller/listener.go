package controller

import (
	"mals/internal/client"
	"mals/pkg/config"
)

type ListenerController interface {
	Shutdown() error
	Serve(onReady func()) error

	Register(config config.Listener) error
	Unregister(name string) error
	Create(name string) error
	Delete(name string) error
	Start(name string) error
	Stop(name string) error

	ClientAdd(name string, client client.Client) error
	ClientRemove(name string, client client.Client) error
}
