package logging

import (
	"fmt"
	"mals/internal/log"
	"mals/pkg/config"
)

func (s *LogController) Register(config config.Log) error {
	if s.internal != nil {
		e := EventRegister{Config: config}
		s.internal <- &e
		return e.Error
	}
	return fmt.Errorf("%T done", s)
}

func (s *LogController) Unregister(name string) error {
	if s.internal != nil {
		e := EventUnregister{Name: name}
		s.internal <- &e
		return e.Error
	}
	return fmt.Errorf("%T done", s)
}

func (s *LogController) Create(name string) error {
	if s.internal != nil {
		e := EventCreate{Name: name}
		s.internal <- &e
		return e.Error
	}
	return fmt.Errorf("%T done", s)
}

func (s *LogController) Delete(name string) error {
	if s.internal != nil {
		e := EventDelete{Name: name}
		s.internal <- &e
		return e.Error
	}
	return fmt.Errorf("%T done", s)
}

func (s *LogController) Start(name string) error {
	if s.internal != nil {
		e := EventStart{Name: name}
		s.internal <- &e
		return e.Error
	}
	return fmt.Errorf("%T done", s)
}

func (s *LogController) Stop(name string) error {
	if s.internal != nil {
		e := EventStop{Name: name}
		s.internal <- &e
		return e.Error
	}
	return fmt.Errorf("%T done", s)
}

func (s *LogController) Debugf(format string, a ...any) error {
	if s.internal != nil {
		e := EventLog{Level: log.LevelDebug, Msg: fmt.Sprintf(format, a...)}
		s.internal <- &e
		return e.Error
	}
	return fmt.Errorf("%T done", s)
}

func (s *LogController) Infof(format string, a ...any) error {
	if s.internal != nil {
		e := EventLog{Level: log.LevelInfo, Msg: fmt.Sprintf(format, a...)}
		s.internal <- &e
		return e.Error
	}
	return fmt.Errorf("%T done", s)
}

func (s *LogController) Warnf(format string, a ...any) error {
	if s.internal != nil {
		e := EventLog{Level: log.LevelWarn, Msg: fmt.Sprintf(format, a...)}
		s.internal <- &e
		return e.Error
	}
	return fmt.Errorf("%T done", s)
}

func (s *LogController) Errorf(format string, a ...any) error {
	if s.internal != nil {
		e := EventLog{Level: log.LevelError, Msg: fmt.Sprintf(format, a...)}
		s.internal <- &e
		return e.Error
	}
	return fmt.Errorf("%T done", s)
}
