package client

import (
	"mals/internal/client"
	"mals/internal/listener"
)

func (s *ClientController) Shutdown() error {
	e := TaskShutdown{TaskGeneric: NewTaskSingle()}
	s.internal <- &e
	return <-e.Result
}

func (s *ClientController) Terminate() error {
	e := TaskTerminate{TaskGeneric: NewTaskSingle()}
	s.internal <- &e
	return <-e.Result
}

func (s *ClientController) Own(client client.Client, listener listener.Listener) error {
	e := TaskOwn{TaskGeneric: NewTaskSingle(), Client: client, Listener: listener.Name()}
	s.internal <- &e
	return <-e.Result
}

func (s *ClientController) Delete(client client.Client) error {
	e := TaskDelete{TaskGeneric: NewTaskSingle(), Client: client}
	s.internal <- &e
	return <-e.Result
}

func (s *ClientController) Start(client client.Client) error {
	e := TaskStart{TaskGeneric: NewTaskSingle(), Client: client}
	s.internal <- &e
	return <-e.Result
}

func (s *ClientController) Stop(client client.Client) error {
	e := TaskStop{TaskGeneric: NewTaskSingle(), Client: client}
	s.internal <- &e
	return <-e.Result
}
