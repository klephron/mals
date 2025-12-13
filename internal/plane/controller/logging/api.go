package logging

import (
	"fmt"
	"mals/internal/log"
	"mals/internal/log/file"
	"mals/pkg/config"
	"sync"
)

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

func (s *LogController) Register(logConfig config.Log) error {
	name := logConfig.Name()

	if _, ok := s.state.logs.Load(name); ok {
		return fmt.Errorf("log %v exists", name)
	}

	switch config := logConfig.(type) {
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
	value, ok := s.state.logs.Load(name)
	if !ok {
		return fmt.Errorf("log %v does not exist", name)
	}

	value.rw.RLock()

	if value.log != nil {
		value.rw.RUnlock()
		return fmt.Errorf("log %v is already created", name)
	}

	s.state.logs.Delete(name)
	value.rw.RUnlock()

	return nil
}

func (s *LogController) Create(name string) error {
	value, ok := s.state.logs.Load(name)
	if !ok {
		return fmt.Errorf("log %v does not exist", name)
	}

	value.rw.Lock()
	defer value.rw.Unlock()

	if value.log != nil {
		return fmt.Errorf("log %v is already created", name)
	}

	switch config := value.config.(type) {
	case *config.LogFile:
		log, err := file.Open(config.File, config.Level)
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
	value, ok := s.state.logs.Load(name)
	if !ok {
		return fmt.Errorf("log %v does not exist", name)
	}

	value.rw.Lock()

	if value.log == nil {
		value.rw.Unlock()
		return fmt.Errorf("log %v is not created", name)
	}
	if value.enabled {
		value.rw.Unlock()
		return fmt.Errorf("log %v is running", name)
	}

	log := value.log
	value.log = nil
	value.rw.Unlock()

	return log.Close()
}

func (s *LogController) Start(name string) error {
	value, ok := s.state.logs.Load(name)
	if !ok {
		return fmt.Errorf("log %v does not exist", name)
	}

	value.rw.Lock()
	if value.log == nil {
		value.rw.Unlock()
		return fmt.Errorf("log %v is not created", name)
	}

	value.enabled = true
	value.rw.Unlock()

	return nil
}

func (s *LogController) Stop(name string) error {
	value, ok := s.state.logs.Load(name)
	if !ok {
		return fmt.Errorf("log %v does not exist", name)
	}

	value.rw.Lock()
	if value.log == nil {
		value.rw.Unlock()
		return fmt.Errorf("log %v is not created", name)
	}

	value.enabled = false
	value.rw.Unlock()

	return nil
}

func (s *LogController) log(level log.Level, format string, a ...any) error {
	s.state.logs.Range(func(key string, value *LogValue) bool {
		value.rw.RLock()
		defer value.rw.RUnlock()

		if value.log == nil || !value.enabled {
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
