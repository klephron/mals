package controller

import "mals/pkg/config"

type LogController interface {
	Serve(onReady func()) error

	Register(config config.Log) error
	Unregister(name string) error
	Create(name string) error
	Delete(name string) error
	Start(name string) error
	Stop(name string) error

	Debugf(format string, a ...any) error
	Infof(format string, a ...any) error
	Warnf(format string, a ...any) error
	Errorf(format string, a ...any) error
}
