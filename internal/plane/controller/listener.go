package controller

import (
	"mals/internal/listener"
	"mals/pkg/config"
)

type ListenerStatus int32

const (
	ListenerAbsent     ListenerStatus = 0
	ListenerRegistered ListenerStatus = (1 << 0)
	ListenerCreated    ListenerStatus = (1 << 1)
	ListenerStarted    ListenerStatus = (1 << 2)
)

type ClientStatus int32

const (
	ClientAbsent  ClientStatus = 0
	ClientCreated ClientStatus = (1 << 1)
	ClientStarted ClientStatus = (1 << 2)
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

	ListenerLspClientOwn(name string, client listener.ListenerLspClient) error
	ListenerLspClientStatus(name string, clientName string) ClientStatus
	ListenerLspClientServe(name string, clientName string) error
	ListenerLspClientShutdown(name string, clientName string) error

	ListenerGetConfig(name string) (*config.Listener, error)
	ListenerGet(name string) (*ListenerData, error)
	ListenerGetAll() []*ListenerData
}
