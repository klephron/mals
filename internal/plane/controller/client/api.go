package client

import (
	"mals/internal/client"
	"mals/internal/listener"
)

func (s *ClientController) Shutdown() error {
	e := TaskShutdown{TaskGeneric: NewTask()}
	s.internal <- &e
	return <-e.Result
}

func (s *ClientController) Own(client client.Client, listener listener.Listener) error {
	e := TaskOwn{TaskGeneric: NewTask(), Client: client, Listener: listener.Name()}
	s.internal <- &e
	return <-e.Result
}

func (s *ClientController) Delete(client client.Client) error {
	e := TaskDelete{TaskGeneric: NewTask(), Client: client, Notify: true}
	s.internal <- &e
	return <-e.Result
}

func (s *ClientController) DeleteSilent(client client.Client) error {
	e := TaskDelete{TaskGeneric: NewTask(), Client: client, Notify: false}
	s.internal <- &e
	return <-e.Result
}

func (s *ClientController) Start(client client.Client) error {
	e := TaskStart{TaskGeneric: NewTask(), Client: client}
	s.internal <- &e
	return <-e.Result
}

func (s *ClientController) Stop(client client.Client) error {
	e := TaskStop{TaskGeneric: NewTask(), Client: client}
	s.internal <- &e
	return <-e.Result
}
