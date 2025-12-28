package controller

import (
	"mals/internal/client"
	"mals/pkg/config"
)

type ListenerStatus int32

const (
	ListenerAbsent     ListenerStatus = 0
	ListenerRegistered ListenerStatus = (1 << 0)
	ListenerCreated    ListenerStatus = (1 << 1)
	ListenerStarted    ListenerStatus = (1 << 2)
)

type ListenerController interface {
	Shutdown() error
	Serve(onReady func()) error

	Status(name string) ListenerStatus
	Register(name string, config config.Listener) error
	Unregister(name string) error
	Create(name string) error
	Delete(name string) error
	Start(name string) error
	Stop(name string) error

	ClientAdd(name string, client client.Client) error
	ClientRemove(name string, client client.Client) error
}
