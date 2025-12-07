package controller

import "mals/pkg/config"

type ListenerController interface {
	Serve(onReady func()) error

	Register(config config.Listener) error
	Unregister(name string) error
	Create(name string) error
	Delete(name string) error
	Start(name string) error
	Stop(name string) error
}
