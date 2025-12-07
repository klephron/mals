package logging

import (
	"fmt"
	"mals/internal/log"
	"mals/pkg/config"
)

func (s *LogController) Shutdown() error {
	e := TaskShutdown{TaskGeneric: NewTaskSingle()}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Register(config config.Log) error {
	e := TaskRegister{TaskGeneric: NewTaskSingle(), Config: config}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Unregister(name string) error {
	e := TaskUnregister{TaskGeneric: NewTaskSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Create(name string) error {
	e := TaskCreate{TaskGeneric: NewTaskSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Delete(name string) error {
	e := TaskDelete{TaskGeneric: NewTaskSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Start(name string) error {
	e := TaskStart{TaskGeneric: NewTaskSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Stop(name string) error {
	e := TaskStop{TaskGeneric: NewTaskSingle(), Name: name}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Debugf(format string, a ...any) error {
	e := TaskLog{TaskGeneric: NewTaskSingle(), Level: log.LevelDebug, Msg: fmt.Sprintf(format, a...)}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Infof(format string, a ...any) error {
	e := TaskLog{TaskGeneric: NewTaskSingle(), Level: log.LevelInfo, Msg: fmt.Sprintf(format, a...)}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Warnf(format string, a ...any) error {
	e := TaskLog{TaskGeneric: NewTaskSingle(), Level: log.LevelWarn, Msg: fmt.Sprintf(format, a...)}
	s.internal <- &e
	return <-e.Result
}

func (s *LogController) Errorf(format string, a ...any) error {
	e := TaskLog{TaskGeneric: NewTaskSingle(), Level: log.LevelError, Msg: fmt.Sprintf(format, a...)}
	s.internal <- &e
	return <-e.Result
}
