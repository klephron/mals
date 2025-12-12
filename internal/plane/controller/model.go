package controller

import "mals/pkg/config"

type ModelController interface {
	Shutdown() error
	Serve(onReady func()) error

	Register(config config.Model) error
	Unregister(name string) error
	Create(name string) error
	Delete(name string) error
	Start(name string) error
	Stop(name string) error
}
