package listener

import (
	"fmt"
	"mals/pkg/config"
)

func (s *ListenerController) Register(config config.Listener) error {
	if s.internal != nil {
		e := EventRegister{Config: config}
		s.internal <- &e
		return e.Error
	}
	return fmt.Errorf("%T done", s)
}

func (s *ListenerController) Unregister(name string) error {
	if s.internal != nil {
		e := EventUnregister{Name: name}
		s.internal <- &e
		return e.Error
	}
	return fmt.Errorf("%T done", s)
}

func (s *ListenerController) Create(name string) error {
	if s.internal != nil {
		e := EventCreate{Name: name}
		s.internal <- &e
		return e.Error
	}
	return fmt.Errorf("%T done", s)
}

func (s *ListenerController) Delete(name string) error {
	if s.internal != nil {
		e := EventDelete{Name: name}
		s.internal <- &e
		return e.Error
	}
	return fmt.Errorf("%T done", s)
}

func (s *ListenerController) Start(name string) error {
	if s.internal != nil {
		e := EventStart{Name: name}
		s.internal <- &e
		return e.Error
	}
	return fmt.Errorf("%T done", s)
}

func (s *ListenerController) Stop(name string) error {
	if s.internal != nil {
		e := EventStop{Name: name}
		s.internal <- &e
		return e.Error
	}
	return fmt.Errorf("%T done", s)
}
