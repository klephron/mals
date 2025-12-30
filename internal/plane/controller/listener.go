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
	Run(onReady func()) error
	Shutdown() error

	ListenerStatus(name string) ListenerStatus
	ListenerRegister(name string, config config.Listener) error
	ListenerUnregister(name string) error
	ListenerCreate(name string) error
	ListenerDelete(name string) error
	ListenerStart(name string) error
	ListenerStop(name string) error

	ListenerClientAdd(name string, client client.Client) error
	ListenerClientRemove(name string, client client.Client) error
}
