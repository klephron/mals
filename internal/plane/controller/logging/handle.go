package logging

import (
	"fmt"
	"mals/internal/log"
	"mals/internal/log/file"
	"mals/pkg/config"
)

func (s *LogController) handleShutdown(t *TaskShutdown) {
	defer close(t.result)

	s.state.logs.Range(func(key string, value *LogValue) bool {
		ts := &TaskStop{TaskGeneric: newTask(), name: key}
		s.handleStop(ts)
		<-ts.result

		td := &TaskDelete{TaskGeneric: newTask(), name: key}
		s.handleDelete(td)
		<-td.result
		return true
	})

	t.result <- nil
}

func (s *LogController) handleLog(t *TaskLog) {
	defer close(t.result)
	s.state.logs.Range(func(key string, value *LogValue) bool {
		if value.log == nil || !value.enabled {
			return true
		}
		switch t.level {
		case log.LevelDebug:
			value.log.Debug(t.msg)
		case log.LevelInfo:
			value.log.Info(t.msg)
		case log.LevelWarn:
			value.log.Warn(t.msg)
		case log.LevelError:
			value.log.Error(t.msg)
		}
		return true
	})

	t.result <- nil
}

func (s *LogController) handleRegister(t *TaskRegister) {
	defer close(t.result)
	name := t.config.Name()

	if _, ok := s.state.logs.Load(name); ok {
		t.result <- fmt.Errorf("log %v exists", name)
		return
	}

	switch config := t.config.(type) {
	case *config.LogFile:
		s.state.logs.Store(name, &LogValue{
			config:  config,
			log:     nil,
			enabled: false,
		})
		t.result <- nil

	default:
		t.result <- fmt.Errorf("unhandled log %T %v", config, config)
	}
}

func (s *LogController) handleUnregister(t *TaskUnregister) {
	defer close(t.result)
	name := t.name

	value, ok := s.state.logs.Load(name)
	if !ok {
		t.result <- fmt.Errorf("log %v does not exist", name)
		return
	}
	if value.log != nil {
		t.result <- fmt.Errorf("log %v is already created", name)
		return
	}

	s.state.logs.Delete(name)

	t.result <- nil
}

func (s *LogController) handleCreate(t *TaskCreate) {
	defer close(t.result)
	name := t.name

	value, ok := s.state.logs.Load(name)
	if !ok {
		t.result <- fmt.Errorf("log %v does not exist", name)
		return
	}
	if value.log != nil {
		t.result <- fmt.Errorf("log %v is already created", name)
		return
	}

	switch config := value.config.(type) {
	case *config.LogFile:
		log, err := file.Open(config.File, config.Level)
		if err != nil {
			t.result <- err
			return
		}
		value.log = log
		value.enabled = false
		t.result <- nil
	default:
		t.result <- fmt.Errorf("unhandled log %T %v", config, config)
	}
}

func (s *LogController) handleDelete(t *TaskDelete) {
	defer close(t.result)
	name := t.name

	value, ok := s.state.logs.Load(name)
	if !ok {
		t.result <- fmt.Errorf("log %v does not exist", name)
		return
	}
	if value.log == nil {
		t.result <- fmt.Errorf("log %v is not created", name)
		return
	}
	if value.enabled {
		t.result <- fmt.Errorf("log %v is running", name)
		return
	}

	t.result <- value.log.Close()
	value.log = nil
}

func (s *LogController) handleStart(t *TaskStart) {
	defer close(t.result)
	name := t.name

	value, ok := s.state.logs.Load(name)
	if !ok {
		t.result <- fmt.Errorf("log %v does not exist", name)
		return
	}
	if value.log == nil {
		t.result <- fmt.Errorf("log %v is not created", name)
		return
	}

	value.enabled = true

	t.result <- nil
}

func (s *LogController) handleStop(t *TaskStop) {
	defer close(t.result)
	name := t.name

	value, ok := s.state.logs.Load(name)
	if !ok {
		t.result <- fmt.Errorf("log %v does not exist", name)
		return
	}
	if value.log == nil {
		t.result <- fmt.Errorf("log %v is not created", name)
		return
	}

	value.enabled = false

	t.result <- nil
}
