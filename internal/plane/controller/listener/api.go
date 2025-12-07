package listener

import "mals/pkg/config"

func (s *ListenerController) Register(config config.Listener) error {
	e := EventRegister{Config: config}
	s.internal <- &e
	return e.Error
}

func (s *ListenerController) Unregister(name string) error {
	e := EventUnregister{Name: name}
	s.internal <- &e
	return e.Error
}

func (s *ListenerController) Create(name string) error {
	e := EventCreate{Name: name}
	s.internal <- &e
	return e.Error
}

func (s *ListenerController) Delete(name string) error {
	e := EventDelete{Name: name}
	s.internal <- &e
	return e.Error
}

func (s *ListenerController) Start(name string) error {
	e := EventStart{Name: name}
	s.internal <- &e
	return e.Error
}

func (s *ListenerController) Stop(name string) error {
	e := EventStop{Name: name}
	s.internal <- &e
	return e.Error
}
