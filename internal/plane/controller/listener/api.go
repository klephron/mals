package listener

import (
	"fmt"
	"mals/pkg/config"
)

func (s *ListenerController) Register(config config.Listener) error {
	if s.internal != nil {
		e := EventRegister{EventGeneric: NewEventSingle(), Config: config}
		s.internal <- &e
		return <-e.Result
	}
	return fmt.Errorf("%T done", s)
}

func (s *ListenerController) Unregister(name string) error {
	if s.internal != nil {
		e := EventUnregister{EventGeneric: NewEventSingle(), Name: name}
		s.internal <- &e
		return <-e.Result
	}
	return fmt.Errorf("%T done", s)
}

func (s *ListenerController) Create(name string) error {
	if s.internal != nil {
		e := EventCreate{EventGeneric: NewEventSingle(), Name: name}
		s.internal <- &e
		return <-e.Result
	}
	return fmt.Errorf("%T done", s)
}

func (s *ListenerController) Delete(name string) error {
	if s.internal != nil {
		e := EventDelete{EventGeneric: NewEventSingle(), Name: name}
		s.internal <- &e
		return <-e.Result
	}
	return fmt.Errorf("%T done", s)
}

func (s *ListenerController) Start(name string) error {
	if s.internal != nil {
		e := EventStart{EventGeneric: NewEventSingle(), Name: name}
		s.internal <- &e
		return <-e.Result
	}
	return fmt.Errorf("%T done", s)
}

func (s *ListenerController) Stop(name string) error {
	if s.internal != nil {
		e := EventStop{EventGeneric: NewEventSingle(), Name: name}
		s.internal <- &e
		return <-e.Result
	}
	return fmt.Errorf("%T done", s)
}
