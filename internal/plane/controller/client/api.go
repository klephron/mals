package client

import (
	"mals/internal/client"
	"mals/internal/listener"
)

func (s *ClientController) Shutdown() error {
	e := TaskShutdown{TaskGeneric: newTask()}
	s.state.internal <- &e
	return <-e.result
}

func (s *ClientController) Own(client client.Client, listener listener.Listener) error {
	e := TaskOwn{TaskGeneric: newTask(), client: client, listener: listener.Name()}
	s.state.internal <- &e
	return <-e.result
}

func (s *ClientController) Delete(client client.Client) error {
	e := TaskDelete{TaskGeneric: newTask(), client: client, notify: true}
	s.state.internal <- &e
	return <-e.result
}

func (s *ClientController) DeleteSilent(client client.Client) error {
	e := TaskDelete{TaskGeneric: newTask(), client: client, notify: false}
	s.state.internal <- &e
	return <-e.result
}

func (s *ClientController) Start(client client.Client) error {
	e := TaskStart{TaskGeneric: newTask(), client: client}
	s.state.internal <- &e
	return <-e.result
}

func (s *ClientController) Stop(client client.Client) error {
	e := TaskStop{TaskGeneric: newTask(), client: client}
	s.state.internal <- &e
	return <-e.result
}
