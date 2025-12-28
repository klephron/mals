package controller

import (
	"mals/pkg/config"
)

type UsageController interface {
	Shutdown() error
	Serve(onReady func()) error

	Register(config config.Usage) error
	Unregister(name string) error

	GetAll() []*config.Usage
}
