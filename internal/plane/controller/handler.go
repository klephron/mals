package controller

import (
	"mals/pkg/config"
)

type HandlerController interface {
	ControllerRun(onReady func()) error
	ControllerShutdown() error

	Register(config *config.Handler) error
	Unregister(name string) error
	Get(name string) (*config.Handler, error)
	GetAll() []*config.Handler
}
