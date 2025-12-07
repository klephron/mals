package logging

import (
	"fmt"
	"mals/internal/log"
	"mals/pkg/config"
)

func (s *LogController) Register(config config.Log) error {
	e := EventRegister{EventGeneric: NewEventSingle(), Config: config}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Unregister(name string) error {
	e := EventUnregister{EventGeneric: NewEventSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Create(name string) error {
	e := EventCreate{EventGeneric: NewEventSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Delete(name string) error {
	e := EventDelete{EventGeneric: NewEventSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Start(name string) error {
	e := EventStart{EventGeneric: NewEventSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Stop(name string) error {
	e := EventStop{EventGeneric: NewEventSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Debugf(format string, a ...any) error {
	e := EventLog{EventGeneric: NewEventSingle(), Level: log.LevelDebug, Msg: fmt.Sprintf(format, a...)}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Infof(format string, a ...any) error {
	e := EventLog{EventGeneric: NewEventSingle(), Level: log.LevelInfo, Msg: fmt.Sprintf(format, a...)}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Warnf(format string, a ...any) error {
	e := EventLog{EventGeneric: NewEventSingle(), Level: log.LevelWarn, Msg: fmt.Sprintf(format, a...)}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Errorf(format string, a ...any) error {
	e := EventLog{EventGeneric: NewEventSingle(), Level: log.LevelError, Msg: fmt.Sprintf(format, a...)}
	s.internal <- &e
	return <-e.Result
}
