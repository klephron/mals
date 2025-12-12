package logging

import (
	"fmt"
	"mals/internal/log"
	"mals/pkg/config"
)

func (s *LogController) Shutdown() error {
	e := TaskShutdown{TaskGeneric: NewTask()}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Register(config config.Log) error {
	e := TaskRegister{TaskGeneric: NewTask(), Config: config}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Unregister(name string) error {
	e := TaskUnregister{TaskGeneric: NewTask(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Create(name string) error {
	e := TaskCreate{TaskGeneric: NewTask(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Delete(name string) error {
	e := TaskDelete{TaskGeneric: NewTask(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Start(name string) error {
	e := TaskStart{TaskGeneric: NewTask(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Stop(name string) error {
	e := TaskStop{TaskGeneric: NewTask(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Debugf(format string, a ...any) error {
	e := TaskLog{TaskGeneric: NewTask(), Level: log.LevelDebug, Msg: fmt.Sprintf(format, a...)}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Infof(format string, a ...any) error {
	e := TaskLog{TaskGeneric: NewTask(), Level: log.LevelInfo, Msg: fmt.Sprintf(format, a...)}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Warnf(format string, a ...any) error {
	e := TaskLog{TaskGeneric: NewTask(), Level: log.LevelWarn, Msg: fmt.Sprintf(format, a...)}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Errorf(format string, a ...any) error {
	e := TaskLog{TaskGeneric: NewTask(), Level: log.LevelError, Msg: fmt.Sprintf(format, a...)}
	s.internal <- &e
	return <-e.Result
}
