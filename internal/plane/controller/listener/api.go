package listener

import (
	"mals/internal/client"
	"mals/pkg/config"
)

func (s *ListenerController) Shutdown() error {
	e := TaskShutdown{TaskGeneric: newTask()}
	s.state.internal <- &e
	return <-e.result
}

func (s *ListenerController) Register(config config.Listener) error {
	e := TaskRegister{TaskGeneric: newTask(), config: config}
	s.state.internal <- &e
	return <-e.result
}

func (s *ListenerController) Unregister(name string) error {
	e := TaskUnregister{TaskGeneric: newTask(), name: name}
	s.state.internal <- &e
	return <-e.result
}

func (s *ListenerController) Create(name string) error {
	e := TaskCreate{TaskGeneric: newTask(), name: name}
	s.state.internal <- &e
	return <-e.result
}

func (s *ListenerController) Delete(name string) error {
	e := TaskDelete{TaskGeneric: newTask(), name: name}
	s.state.internal <- &e
	return <-e.result
}

func (s *ListenerController) Start(name string) error {
	e := TaskStart{TaskGeneric: newTask(), name: name}
	s.state.internal <- &e
	return <-e.result
}

func (s *ListenerController) Stop(name string) error {
	e := TaskStop{TaskGeneric: newTask(), name: name}
	s.state.internal <- &e
	return <-e.result
}

func (s *ListenerController) ClientAdd(name string, client client.Client) error {
	e := TaskClientAdd{TaskGeneric: newTask(), name: name, client: client}
	s.state.internal <- &e
	return <-e.result
}

func (s *ListenerController) ClientRemove(name string, client client.Client) error {
	e := TaskClientRemove{TaskGeneric: newTask(), name: name, client: client}
	s.state.internal <- &e
	return <-e.result
}
