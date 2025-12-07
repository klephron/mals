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
		s.state.Logs.Delete(key)
		return true
	})
}

func (s *LogController) handleLog(e *EventLog) {
	defer close(e.Result)
	s.state.Logs.Range(func(key string, value *state.LogValue) bool {
		if value.Log == nil || !value.Enabled {
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
	e.Result <- nil
}

func (s *LogController) handleRegister(e *EventRegister) {
	defer close(e.Result)
	name := e.Config.Name()

	if _, ok := s.state.Logs.Load(name); ok {
		e.Result <- fmt.Errorf("log %v exists", name)
		return
	}

	switch config := e.Config.(type) {
	case *config.LogFile:
		s.state.Logs.Store(name, &state.LogValue{
			Config:  config,
			Log:     nil,
			Enabled: false,
		})
		e.Result <- nil

	default:
		e.Result <- fmt.Errorf("unhandled log %T %v", config, config)
	}
}

func (s *LogController) handleUnregister(e *EventUnregister) {
	defer close(e.Result)
	name := e.Name

	value, ok := s.state.Logs.Load(name)
	if !ok {
		e.Result <- fmt.Errorf("log %v does not exist", name)
		return
	}

	if value.Log != nil {
		e.Result <- fmt.Errorf("log %v is already created", name)
		return
	}

	s.state.Logs.Delete(name)
	e.Result <- nil
}

func (s *LogController) handleCreate(e *EventCreate) {
	defer close(e.Result)
	name := e.Name

	value, ok := s.state.Logs.Load(name)
	if !ok {
		e.Result <- fmt.Errorf("log %v does not exist", name)
		return
	}

	if value.Log != nil {
		e.Result <- fmt.Errorf("log %v is already created", name)
		return
	}

	switch config := value.Config.(type) {
	case *config.LogFile:
		log, err := file.Open(config.File, config.Level)
		if err != nil {
			e.Result <- err
			return
		}
		value.Log = log
		value.Enabled = false
		e.Result <- nil
	default:
		e.Result <- fmt.Errorf("unhandled log %T %v", config, config)
	}
}

func (s *LogController) handleDelete(e *EventDelete) {
	defer close(e.Result)
	name := e.Name

	value, ok := s.state.Logs.Load(name)
	if !ok {
		e.Result <- fmt.Errorf("log %v does not exist", name)
		return
	}

	if value.Log == nil {
		e.Result <- fmt.Errorf("log %v is not created", name)
		return
	}

	if value.Enabled {
		e.Result <- fmt.Errorf("log %v is running", name)
		return
	}

	e.Result <- value.Log.Close()
	value.Log = nil
}

func (s *LogController) handleStart(e *EventStart) {
	defer close(e.Result)
	name := e.Name

	value, ok := s.state.Logs.Load(name)
	if !ok {
		e.Result <- fmt.Errorf("log %v does not exist", name)
		return
	}

	if value.Log == nil {
		e.Result <- fmt.Errorf("log %v is not created", name)
		return
	}

	value.Enabled = true
	e.Result <- nil
}

func (s *LogController) handleStop(e *EventStop) {
	defer close(e.Result)
	name := e.Name

	value, ok := s.state.Logs.Load(name)
	if !ok {
		e.Result <- fmt.Errorf("log %v does not exist", name)
		return
	}

	if value.Log == nil {
		e.Result <- fmt.Errorf("log %v is not created", name)
		return
	}

	value.Enabled = false
	e.Result <- nil
}
