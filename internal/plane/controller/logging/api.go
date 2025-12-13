package logging

import (
	"fmt"
	"mals/internal/log"
	"mals/pkg/config"
)

func (s *LogController) Shutdown() error {
	e := TaskShutdown{TaskGeneric: newTask()}
	s.state.taskChan <- &e
	return <-e.result
}

func (s *LogController) Register(config config.Log) error {
	e := TaskRegister{TaskGeneric: newTask(), config: config}
	s.state.taskChan <- &e
	return <-e.result
}

func (s *LogController) Unregister(name string) error {
	e := TaskUnregister{TaskGeneric: newTask(), name: name}
	s.state.taskChan <- &e
	return <-e.result
}

func (s *LogController) Create(name string) error {
	e := TaskCreate{TaskGeneric: newTask(), name: name}
	s.state.taskChan <- &e
	return <-e.result
}

func (s *LogController) Delete(name string) error {
	e := TaskDelete{TaskGeneric: newTask(), name: name}
	s.state.taskChan <- &e
	return <-e.result
}

func (s *LogController) Start(name string) error {
	e := TaskStart{TaskGeneric: newTask(), name: name}
	s.state.taskChan <- &e
	return <-e.result
}

func (s *LogController) Stop(name string) error {
	e := TaskStop{TaskGeneric: newTask(), name: name}
	s.state.taskChan <- &e
	return <-e.result
}

func (s *LogController) Debugf(format string, a ...any) error {
	e := TaskLog{TaskGeneric: newTask(), level: log.LevelDebug, msg: fmt.Sprintf(format, a...)}
	s.state.taskChan <- &e
	return <-e.result
}

func (s *LogController) Infof(format string, a ...any) error {
	e := TaskLog{TaskGeneric: newTask(), level: log.LevelInfo, msg: fmt.Sprintf(format, a...)}
	s.state.taskChan <- &e
	return <-e.result
}

func (s *LogController) Warnf(format string, a ...any) error {
	e := TaskLog{TaskGeneric: newTask(), level: log.LevelWarn, msg: fmt.Sprintf(format, a...)}
	s.state.taskChan <- &e
	return <-e.result
}

func (s *LogController) Errorf(format string, a ...any) error {
	e := TaskLog{TaskGeneric: newTask(), level: log.LevelError, msg: fmt.Sprintf(format, a...)}
	s.state.taskChan <- &e
	return <-e.result
}
