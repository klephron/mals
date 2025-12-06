package logging

import (
	"fmt"
	"mals/internal/log"
	"mals/internal/log/file"
	"mals/internal/plane/event"
	"mals/internal/plane/state"
	"mals/pkg/config"
)

func (s *LogController) handleShutdown(_ *event.EventShutdown) {
	s.state.Logs.Range(func(key string, value *state.LogValue) bool {
		s.handleStop(&EventStop{Name: key})
		s.handleDelete(&EventDelete{Name: key})
		return true
	})
}

func (s *LogController) handleTerminate(_ *event.EventTerminate) {
	s.state.Logs.Range(func(key string, value *state.LogValue) bool {
		s.handleDelete(&EventDelete{Name: key})
		return true
	})
}

func (s *LogController) handleLog(e *event.EventLog) {
	s.state.Logs.Range(func(key string, value *state.LogValue) bool {
		if value.Log == nil && value.Enabled {
			return true
		}
		switch e.Level {
		case log.LevelDebug:
			value.Log.Debug(e.Msg)
		case log.LevelInfo:
			value.Log.Info(e.Msg)
		case log.LevelWarn:
			value.Log.Warn(e.Msg)
		case log.LevelError:
			value.Log.Error(e.Msg)
		}
		return true
	})
}

func (s *LogController) handleRegister(e *EventRegister) {
	name := e.Config.Name()

	if _, ok := s.state.Logs.Load(name); ok {
		e.Error = fmt.Errorf("log %v exists", name)
		return
	}

	switch config := e.Config.(type) {
	case *config.LogFile:
		s.state.Logs.Store(name, &state.LogValue{
			Config:  config,
			Log:     nil,
			Enabled: false,
		})
		e.Error = nil

	default:
		e.Error = fmt.Errorf("unhandled log %T %v", config, config)
	}
}

func (s *LogController) handleUnregister(e *EventUnregister) {
	name := e.Name

	value, ok := s.state.Logs.Load(name)
	if !ok {
		e.Error = fmt.Errorf("log %v does not exist", name)
		return
	}

	if value.Log != nil {
		e.Error = fmt.Errorf("log %v is already created", name)
		return
	}

	s.state.Logs.Delete(name)
	e.Error = nil
}

func (s *LogController) handleCreate(e *EventCreate) {
	name := e.Name

	value, ok := s.state.Logs.Load(name)
	if !ok {
		e.Error = fmt.Errorf("log %v does not exist", name)
		return
	}

	if value.Log != nil {
		e.Error = fmt.Errorf("log %v is already created", name)
		return
	}

	switch config := value.Config.(type) {
	case *config.LogFile:
		log, err := file.Open(config.File, config.Level)
		if err != nil {
			e.Error = err
			return
		}
		value.Log = log
		value.Enabled = false
		e.Error = nil
	default:
		e.Error = fmt.Errorf("unhandled log %T %v", config, config)
	}
}

func (s *LogController) handleDelete(e *EventDelete) {
	name := e.Name

	value, ok := s.state.Logs.Load(name)
	if !ok {
		e.Error = fmt.Errorf("log %v does not exist", name)
		return
	}

	if value.Log == nil {
		e.Error = fmt.Errorf("log %v is not created", name)
		return
	}

	if value.Enabled {
		e.Error = fmt.Errorf("log %v is running", name)
		return
	}

	e.Error = value.Log.Close()
	value.Log = nil
}

func (s *LogController) handleStart(e *EventStart) {
	name := e.Name

	value, ok := s.state.Logs.Load(name)
	if !ok {
		e.Error = fmt.Errorf("log %v does not exist", name)
		return
	}

	if value.Log == nil {
		e.Error = fmt.Errorf("log %v is not created", name)
		return
	}

	value.Enabled = true
	e.Error = nil
}

func (s *LogController) handleStop(e *EventStop) {
	name := e.Name

	value, ok := s.state.Logs.Load(name)
	if !ok {
		e.Error = fmt.Errorf("log %v does not exist", name)
		return
	}

	if value.Log == nil {
		e.Error = fmt.Errorf("log %v is not created", name)
		return
	}

	value.Enabled = false
	e.Error = nil
}
