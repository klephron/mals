package logging

import (
	"fmt"
	"mals/internal/log"
	"mals/internal/plane/event"
	"mals/pkg/config"
)

func (s *LogController) Register(config config.Log) error {
	e := EventRegister{Config: config}
	s.internal <- &e
	return e.Error
}

func (s *LogController) Unregister(name string) error {
	e := EventUnregister{Name: name}
	s.internal <- &e
	return e.Error
}

func (s *LogController) Create(name string) error {
	e := EventCreate{Name: name}
	s.internal <- &e
	return e.Error
}

func (s *LogController) Delete(name string) error {
	e := EventDelete{Name: name}
	s.internal <- &e
	return e.Error
}

func (s *LogController) Start(name string) error {
	e := EventStart{Name: name}
	s.internal <- &e
	return e.Error
}

func (s *LogController) Stop(name string) error {
	e := EventStop{Name: name}
	s.internal <- &e
	return e.Error
}

func (s *LogController) Debugf(format string, a ...any) {
	s.bus.Unicast(&event.EventLog{
		Level: log.LevelDebug,
		Msg:   fmt.Sprintf(format, a...),
	}, s.external)
}

func (s *LogController) Infof(format string, a ...any) {
	s.bus.Unicast(&event.EventLog{
		Level: log.LevelInfo,
		Msg:   fmt.Sprintf(format, a...),
	}, s.external)
}

func (s *LogController) Warnf(format string, a ...any) {
	s.bus.Unicast(&event.EventLog{
		Level: log.LevelWarn,
		Msg:   fmt.Sprintf(format, a...),
	}, s.external)
}

func (s *LogController) Errorf(format string, a ...any) {
	s.bus.Unicast(&event.EventLog{
		Level: log.LevelError,
		Msg:   fmt.Sprintf(format, a...),
	}, s.external)
}
