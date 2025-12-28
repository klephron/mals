package logging

import (
	"fmt"
	"mals/internal/log"
	"mals/internal/log/file"
	"mals/internal/plane/controller"
	"mals/pkg/config"
	"sync"
)

func statusErrorEq(name string, actual controller.LogStatus, expected controller.LogStatus) error {
	return fmt.Errorf("log %v expected eq %v, got %v", name, expected, actual)
}

func statusErrorFlag(name string, actual controller.LogStatus, expected controller.LogStatus) error {
	return fmt.Errorf("log %v expected flag %v, got %v", name, expected, actual)
}

func (s *LogController) status(value *LogValue) controller.LogStatus {
	status := controller.LogAbsent

	if value != nil {
		status |= controller.LogRegistered

		if value.log != nil {
			status |= controller.LogCreated
		}
		if value.enabled {
			status |= controller.LogStarted
		}
	}

	return status
}

func (s *LogController) statusRW(value *LogValue) controller.LogStatus {
	status := controller.LogAbsent

	if value != nil {
		status |= controller.LogRegistered

		value.rw.RLock()

		if value.log != nil {
			status |= controller.LogCreated
		}
		if value.enabled {
			status |= controller.LogStarted
		}

		value.rw.RUnlock()
	}

	return status
}

func (s *LogController) Shutdown() error {
	s.state.logs.Range(func(key string, value *LogValue) bool {
		s.Stop(key)
		s.Delete(key)
		return true
	})

	s.state.statusRW.RLock()
	cancel := s.state.statusCancel
	s.state.statusRW.RUnlock()

	cancel()

	return nil
}

func (s *LogController) Status(name string) controller.LogStatus {
	value, _ := s.state.logs.Load(name)
	return s.statusRW(value)
}

func (s *LogController) Register(cfg config.Log) error {
	name := cfg.Name()

	value, _ := s.state.logs.Load(name)

	if status := s.statusRW(value); status != controller.LogAbsent {
		return statusErrorEq(name, status, controller.LogAbsent)
	}

	switch config := cfg.(type) {
	case *config.LogFile:
		s.state.logs.Store(name, &LogValue{
			rw:      sync.RWMutex{},
			config:  config,
			log:     nil,
			enabled: false,
		})

	default:
		return fmt.Errorf("unhandled log %T %v", config, config)
	}

	return nil
}

func (s *LogController) Unregister(name string) error {
	value, _ := s.state.logs.Load(name)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status != controller.LogRegistered {
		return statusErrorEq(name, status, controller.LogRegistered)
	}

	s.state.logs.Delete(name)

	return nil
}

func (s *LogController) Create(name string) error {
	value, _ := s.state.logs.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status != controller.LogRegistered {
		return statusErrorEq(name, status, controller.LogRegistered)
	}

	switch config := value.config.(type) {
	case *config.LogFile:
		log, err := file.Open(config.File, config.Level())
		if err != nil {
			return err
		}
		value.log = log
		value.enabled = false
	default:
		return fmt.Errorf("unhandled log %T %v", config, config)
	}

	return nil
}

func (s *LogController) Delete(name string) error {
	value, _ := s.state.logs.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status != controller.LogRegistered|controller.LogCreated {
		return statusErrorEq(name, status, controller.LogRegistered|controller.LogCreated)
	}

	value.log.Close()
	value.log = nil

	return nil
}

func (s *LogController) Start(name string) error {
	value, _ := s.state.logs.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status != controller.LogRegistered|controller.LogCreated {
		return statusErrorEq(name, status, controller.LogRegistered|controller.LogCreated)
	}

	value.enabled = true

	return nil
}

func (s *LogController) Stop(name string) error {
	value, _ := s.state.logs.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status&controller.LogStarted == 0 {
		return statusErrorFlag(name, status, controller.LogStarted)
	}

	value.enabled = false

	return nil
}

func (s *LogController) log(level log.Level, format string, a ...any) error {
	s.state.logs.Range(func(key string, value *LogValue) bool {
		value.rw.RLock()
		defer value.rw.RUnlock()

		if status := s.status(value); status&controller.LogStarted == 0 {
			return true
		}

		msg := fmt.Sprintf(format, a...)

		switch level {
		case log.LevelDebug:
			value.log.Debug(msg)
		case log.LevelInfo:
			value.log.Info(msg)
		case log.LevelWarn:
			value.log.Warn(msg)
		case log.LevelError:
			value.log.Error(msg)
		}
		return true
	})

	return nil
}

func (s *LogController) Debugf(format string, a ...any) error {
	return s.log(log.LevelDebug, format, a...)
}

func (s *LogController) Infof(format string, a ...any) error {
	return s.log(log.LevelInfo, format, a...)
}

func (s *LogController) Warnf(format string, a ...any) error {
	return s.log(log.LevelWarn, format, a...)
}

func (s *LogController) Errorf(format string, a ...any) error {
	return s.log(log.LevelError, format, a...)
}
