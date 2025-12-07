package controller

import (
	"mals/internal/client"
	"mals/internal/listener"
)

type ClientController interface {
	Shutdown() error
	Terminate() error
	Serve(onReady func()) error

	Own(client client.Client, listener listener.Listener) error
	Delete(client client.Client) error
	DeleteSilent(client client.Client) error
	Start(client client.Client) error
	Stop(client client.Client) error
}
