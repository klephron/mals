package logging

import (
	"fmt"
	"mals/internal/log"
	"mals/internal/log/file"
	"mals/internal/plane/state"
	"mals/pkg/config"
)

func (s *LogController) handleLog(t *TaskLog) {
	defer close(t.Result)
	s.state.Logs.Range(func(key string, value *state.LogValue) bool {
		if value.Log == nil || !value.Enabled {
			return true
		}
		switch t.Level {
		case log.LevelDebug:
			value.Log.Debug(t.Msg)
		case log.LevelInfo:
			value.Log.Info(t.Msg)
		case log.LevelWarn:
			value.Log.Warn(t.Msg)
		case log.LevelError:
			value.Log.Error(t.Msg)
		}
		return true
	})
	t.Result <- nil
}

func (s *LogController) handleRegister(t *TaskRegister) {
	defer close(t.Result)
	name := t.Config.Name()

	if _, ok := s.state.Logs.Load(name); ok {
		t.Result <- fmt.Errorf("log %v exists", name)
		return
	}

	switch config := t.Config.(type) {
	case *config.LogFile:
		s.state.Logs.Store(name, &state.LogValue{
			Config:  config,
			Log:     nil,
			Enabled: false,
		})
		t.Result <- nil

	default:
		t.Result <- fmt.Errorf("unhandled log %T %v", config, config)
	}
}

func (s *LogController) handleUnregister(t *TaskUnregister) {
	defer close(t.Result)
	name := t.Name

	value, ok := s.state.Logs.Load(name)
	if !ok {
		t.Result <- fmt.Errorf("log %v does not exist", name)
		return
	}

	if value.Log != nil {
		t.Result <- fmt.Errorf("log %v is already created", name)
		return
	}

	s.state.Logs.Delete(name)
	t.Result <- nil
}

func (s *LogController) handleCreate(t *TaskCreate) {
	defer close(t.Result)
	name := t.Name

	value, ok := s.state.Logs.Load(name)
	if !ok {
		t.Result <- fmt.Errorf("log %v does not exist", name)
		return
	}

	if value.Log != nil {
		t.Result <- fmt.Errorf("log %v is already created", name)
		return
	}

	switch config := value.Config.(type) {
	case *config.LogFile:
		log, err := file.Open(config.File, config.Level)
		if err != nil {
			t.Result <- err
			return
		}
		value.Log = log
		value.Enabled = false
		t.Result <- nil
	default:
		t.Result <- fmt.Errorf("unhandled log %T %v", config, config)
	}
}

func (s *LogController) handleDelete(t *TaskDelete) {
	defer close(t.Result)
	name := t.Name

	value, ok := s.state.Logs.Load(name)
	if !ok {
		t.Result <- fmt.Errorf("log %v does not exist", name)
		return
	}

	if value.Log == nil {
		t.Result <- fmt.Errorf("log %v is not created", name)
		return
	}

	if value.Enabled {
		t.Result <- fmt.Errorf("log %v is running", name)
		return
	}

	t.Result <- value.Log.Close()
	value.Log = nil
}

func (s *LogController) handleStart(t *TaskStart) {
	defer close(t.Result)
	name := t.Name

	value, ok := s.state.Logs.Load(name)
	if !ok {
		t.Result <- fmt.Errorf("log %v does not exist", name)
		return
	}

	if value.Log == nil {
		t.Result <- fmt.Errorf("log %v is not created", name)
		return
	}

	value.Enabled = true
	t.Result <- nil
}

func (s *LogController) handleStop(t *TaskStop) {
	defer close(t.Result)
	name := t.Name

	value, ok := s.state.Logs.Load(name)
	if !ok {
		t.Result <- fmt.Errorf("log %v does not exist", name)
		return
	}

	if value.Log == nil {
		t.Result <- fmt.Errorf("log %v is not created", name)
		return
	}

	value.Enabled = false
	t.Result <- nil
}
