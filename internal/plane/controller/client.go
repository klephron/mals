package controller

import (
	"mals/internal/client"
	"mals/internal/listener"
)

type ClientController interface {
	Shutdown() error
	Run(onReady func()) error

	ClientOwn(client client.Client, listener listener.Listener) error
	ClientServe(client client.Client) error
	ClientShutdown(client client.Client) error
	ClientShutdownSilent(client client.Client) error
}
