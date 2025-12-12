package model

import "mals/pkg/config"

func (s *ModelController) Shutdown() error {
	e := TaskShutdown{TaskGeneric: NewTask()}
	s.internal <- &e
	return <-e.Result
}

func (s *ModelController) Register(config config.Model) error {
	e := TaskRegister{TaskGeneric: NewTask(), Config: config}
	s.internal <- &e
	return <-e.Result
}

func (s *ModelController) Unregister(name string) error {
	e := TaskUnregister{TaskGeneric: NewTask(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *ModelController) Create(name string) error {
	e := TaskCreate{TaskGeneric: NewTask(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *ModelController) Delete(name string) error {
	e := TaskDelete{TaskGeneric: NewTask(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *ModelController) Start(name string) error {
	e := TaskStart{TaskGeneric: NewTask(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *ModelController) Stop(name string) error {
	e := TaskStop{TaskGeneric: NewTask(), Name: name}
	s.internal <- &e
	return <-e.Result
}
