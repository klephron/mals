package controller

import "mals/pkg/config"

type LogController interface {
	Serve() error

	LogRegister(log config.Log) error

	LogCreate(name string)
	LogDelete(name string)
	LogStart(name string)
	LogStop(name string)

	Debugf(format string, a ...any)
	Infof(format string, a ...any)
	Warnf(format string, a ...any)
	Errorf(format string, a ...any)
}
