package controller

import (
	"mals/pkg/config"
)

type ListenerStatus int32

const (
	ListenerAbsent     ListenerStatus = 0
	ListenerRegistered ListenerStatus = (1 << 0)
	ListenerCreated    ListenerStatus = (1 << 1)
	ListenerStarted    ListenerStatus = (1 << 2)
)

type ListenerData struct {
	Name   string
	Status ListenerStatus
	Config *config.Listener
}

type ListenerController interface {
	Run(onReady func()) error
	Shutdown() error

	ListenerStatus(name string) ListenerStatus
	ListenerRegister(name string, config *config.Listener) error
	ListenerUnregister(name string) error
	ListenerCreate(name string) error
	ListenerDelete(name string) error
	ListenerStart(name string) error
	ListenerStop(name string) error

	ListenerClientAdd(name string, client string) error
	ListenerClientRemove(name string, client string) error

	ListenerGetConfig(name string) (*config.Listener, error)
	ListenerGet(name string) (*ListenerData, error)
	ListenerGetAll() []*ListenerData
}
