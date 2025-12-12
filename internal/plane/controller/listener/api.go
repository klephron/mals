package listener

import (
	"mals/internal/client"
	"mals/pkg/config"
)

func (s *ListenerController) Shutdown() error {
	e := TaskShutdown{TaskGeneric: NewTask()}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) Register(config config.Listener) error {
	e := TaskRegister{TaskGeneric: NewTask(), Config: config}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) Unregister(name string) error {
	e := TaskUnregister{TaskGeneric: NewTask(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) Create(name string) error {
	e := TaskCreate{TaskGeneric: NewTask(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) Delete(name string) error {
	e := TaskDelete{TaskGeneric: NewTask(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) Start(name string) error {
	e := TaskStart{TaskGeneric: NewTask(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) Stop(name string) error {
	e := TaskStop{TaskGeneric: NewTask(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) ClientAdd(name string, client client.Client) error {
	e := TaskClientAdd{TaskGeneric: NewTask(), Name: name, Client: client}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) ClientRemove(name string, client client.Client) error {
	e := TaskClientRemove{TaskGeneric: NewTask(), Name: name, Client: client}
	s.internal <- &e
	return <-e.Result
}
