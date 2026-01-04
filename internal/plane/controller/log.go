package controller

import "mals/pkg/config"

type LogStatus int32

const (
	LogAbsent     LogStatus = 0
	LogRegistered LogStatus = (1 << 0)
	LogCreated    LogStatus = (1 << 1)
	LogStarted    LogStatus = (1 << 2)
)

type LogData struct {
	Name   string
	Status LogStatus
	Config *config.Log
}

type LogController interface {
	Run(onReady func()) error
	Shutdown() error

	LogStatus(name string) LogStatus
	LogRegister(name string, config *config.Log) error
	LogUnregister(name string) error
	LogCreate(name string) error
	LogDelete(name string) error
	LogStart(name string) error
	LogStop(name string) error

	LogGet(name string) (*LogData, error)
	LogGetAll() []*LogData

	Debugf(format string, a ...any) error
	Infof(format string, a ...any) error
	Warnf(format string, a ...any) error
	Errorf(format string, a ...any) error
}
