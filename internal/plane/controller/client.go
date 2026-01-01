package controller

import (
	"mals/internal/client"
	"mals/internal/listener"
)

type ClientStatus int32

const (
	ClientAbsent  ClientStatus = 0
	ClientCreated ClientStatus = (1 << 0)
	ClientStarted ClientStatus = (1 << 1)
)

type ClientController interface {
	Shutdown() error
	Run(onReady func()) error

	ClientStatus(name string) ClientStatus
	ClientOwn(name string, client client.Client, listener listener.Listener) error
	ClientServe(name string) error
	ClientShutdown(name string) error
	ClientShutdownSilent(name string) error

	ClientGetListener(name string) (string, error)
}
