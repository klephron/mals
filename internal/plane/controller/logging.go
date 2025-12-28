package controller

import "mals/pkg/config"

type LogStatus int32

const (
	LogAbsent     LogStatus = 0
	LogRegistered LogStatus = (1 << 0)
	LogCreated    LogStatus = (1 << 1)
	LogStarted    LogStatus = (1 << 2)
)

type LogController interface {
	Shutdown() error
	Serve(onReady func()) error

	Status(name string) LogStatus
	Register(name string, config config.Log) error
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
