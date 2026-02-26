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
	ControllerRun(onReady func()) error
	ControllerShutdown() error

	Status(name string) ListenerStatus
	Register(name string, config *config.Listener) error
	Unregister(name string) error
	Create(name string) error
	Delete(name string) error
	Start(name string) error
	Stop(name string) error

	LspClientOwn(name string, client listener.ListenerLspClient) error
	LspClientStatus(name string, clientName string) ClientStatus
	LspClientServe(name string, clientName string) error
	LspClientShutdown(name string, clientName string) error

	GetConfig(name string) (*config.Listener, error)
	Get(name string) (*ListenerData, error)
	GetAll() []*ListenerData
}
