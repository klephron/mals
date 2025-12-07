package listener

import (
	"mals/pkg/config"
)

func (s *ListenerController) Register(config config.Listener) error {
	e := EventRegister{EventGeneric: NewEventSingle(), Config: config}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) Unregister(name string) error {
	e := EventUnregister{EventGeneric: NewEventSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) Create(name string) error {
	e := EventCreate{EventGeneric: NewEventSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) Delete(name string) error {
	e := EventDelete{EventGeneric: NewEventSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) Start(name string) error {
	e := EventStart{EventGeneric: NewEventSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *ListenerController) Stop(name string) error {
	e := EventStop{EventGeneric: NewEventSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}
