package logging

import (
	"fmt"
	"mals/internal/log"
	"mals/pkg/config"
)

func (s *LogController) Register(config config.Log) error {
	if s.internal != nil {
		e := EventRegister{EventGeneric: NewEventSingle(), Config: config}
		s.internal <- &e
		return <-e.Result
	}
	return fmt.Errorf("%T done", s)
}

func (s *LogController) Unregister(name string) error {
	if s.internal != nil {
		e := EventUnregister{EventGeneric: NewEventSingle(), Name: name}
		s.internal <- &e
		return <-e.Result
	}
	return fmt.Errorf("%T done", s)
}

func (s *LogController) Create(name string) error {
	if s.internal != nil {
		e := EventCreate{EventGeneric: NewEventSingle(), Name: name}
		s.internal <- &e
		return <-e.Result
	}
	return fmt.Errorf("%T done", s)
}

func (s *LogController) Delete(name string) error {
	if s.internal != nil {
		e := EventDelete{EventGeneric: NewEventSingle(), Name: name}
		s.internal <- &e
		return <-e.Result
	}
	return fmt.Errorf("%T done", s)
}

func (s *LogController) Start(name string) error {
	if s.internal != nil {
		e := EventStart{EventGeneric: NewEventSingle(), Name: name}
		s.internal <- &e
		return <-e.Result
	}
	return fmt.Errorf("%T done", s)
}

func (s *LogController) Stop(name string) error {
	if s.internal != nil {
		e := EventStop{EventGeneric: NewEventSingle(), Name: name}
		s.internal <- &e
		return <-e.Result
	}
	return fmt.Errorf("%T done", s)
}

func (s *LogController) Debugf(format string, a ...any) error {
	if s.internal != nil {
		e := EventLog{EventGeneric: NewEventSingle(), Level: log.LevelDebug, Msg: fmt.Sprintf(format, a...)}
		s.internal <- &e
		return <-e.Result
	}
	return fmt.Errorf("%T done", s)
}

func (s *LogController) Infof(format string, a ...any) error {
	if s.internal != nil {
		e := EventLog{EventGeneric: NewEventSingle(), Level: log.LevelInfo, Msg: fmt.Sprintf(format, a...)}
		s.internal <- &e
		return <-e.Result
	}
	return fmt.Errorf("%T done", s)
}

func (s *LogController) Warnf(format string, a ...any) error {
	if s.internal != nil {
		e := EventLog{EventGeneric: NewEventSingle(), Level: log.LevelWarn, Msg: fmt.Sprintf(format, a...)}
		s.internal <- &e
		return <-e.Result
	}
	return fmt.Errorf("%T done", s)
}

func (s *LogController) Errorf(format string, a ...any) error {
	if s.internal != nil {
		e := EventLog{EventGeneric: NewEventSingle(), Level: log.LevelError, Msg: fmt.Sprintf(format, a...)}
		s.internal <- &e
		return <-e.Result
	}
	return fmt.Errorf("%T done", s)
}
