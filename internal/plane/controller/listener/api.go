package listener

import (
	"mals/pkg/config"
)

func (s *ListenerController) Shutdown() error {
	e := TaskShutdown{TaskGeneric: NewTaskSingle()}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) Register(config config.Listener) error {
	e := TaskRegister{TaskGeneric: NewTaskSingle(), Config: config}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) Unregister(name string) error {
	e := TaskUnregister{TaskGeneric: NewTaskSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) Create(name string) error {
	e := TaskCreate{TaskGeneric: NewTaskSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) Delete(name string) error {
	e := TaskDelete{TaskGeneric: NewTaskSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) Start(name string) error {
	e := TaskStart{TaskGeneric: NewTaskSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) Stop(name string) error {
	e := TaskStop{TaskGeneric: NewTaskSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}
