package listener

import (
	"context"
	"fmt"
	"mals/internal/listener/api"
	"mals/internal/listener/lsp"
	"mals/internal/plane/event"
	"mals/internal/plane/state"
	"mals/pkg/config"
)

func (s *ListenerController) handleShutdown(_ *event.EventShutdown) {
	s.state.Listeners.Range(func(key string, value *state.ListenerValue) bool {
		s.handleStop(&TaskStop{TaskGeneric: NewTaskSingle(), Name: key})
		s.handleDelete(&TaskDelete{TaskGeneric: NewTaskSingle(), Name: key})
		return true
	})
}

func (s *ListenerController) handleTerminate(_ *event.EventTerminate) {
	s.state.Listeners.Range(func(key string, value *state.ListenerValue) bool {
		s.state.Listeners.Delete(key)
		return true
	})
}

func (s *ListenerController) handleRegister(t *TaskRegister) {
	defer close(t.Result)
	name := t.Config.Name()

	if _, ok := s.state.Listeners.Load(name); ok {
		t.Result <- fmt.Errorf("listener %v exists", name)
		return
	}

	switch config := t.Config.(type) {
	case *config.ListenerTcp:
		s.state.Listeners.Store(name, &state.ListenerValue{
			Config:     config,
			Listener:   nil,
			CancelFunc: nil,
		})
		t.Result <- nil

	default:
		t.Result <- fmt.Errorf("unhandled listener %T %v", config, config)
	}
}

func (s *ListenerController) handleUnregister(t *TaskUnregister) {
	defer close(t.Result)
	name := t.Name

	value, ok := s.state.Listeners.Load(name)
	if !ok {
		t.Result <- fmt.Errorf("listener %v does not exist", name)
		return
	}

	if value.Listener != nil {
		t.Result <- fmt.Errorf("listener %v is already created", name)
		return
	}

	s.state.Listeners.Delete(name)
	t.Result <- nil
}

func (s *ListenerController) handleCreate(t *TaskCreate) {
	defer close(t.Result)
	name := t.Name

	value, ok := s.state.Listeners.Load(name)
	if !ok {
		t.Result <- fmt.Errorf("listener %v does not exist", name)
		return
	}

	if value.Listener != nil {
		t.Result <- fmt.Errorf("listener %v is already created", name)
		return
	}

	switch config := value.Config.(type) {
	case *config.ListenerTcp:
		switch config.Kind() {
		case lsp.Kind():
			if listener, err := lsp.New(name, config.Port, s.plane); err != nil {
				t.Result <- err
			} else {
				value.Listener = listener
				t.Result <- nil
			}

		case api.Kind():
			if listener, err := api.New(name, config.Port, s.plane); err != nil {
				t.Result <- err
			} else {
				value.Listener = listener
				t.Result <- nil
			}

		default:
			t.Result <- fmt.Errorf("unhandled listener %T %v", config, config)
		}

	default:
		t.Result <- fmt.Errorf("unhandled listener %T %v", config, config)
	}
}

func (s *ListenerController) handleDelete(t *TaskDelete) {
	defer close(t.Result)
	name := t.Name

	value, ok := s.state.Listeners.Load(name)
	if !ok {
		t.Result <- fmt.Errorf("listener %v does not exist", name)
		return
	}

	if value.Listener == nil {
		t.Result <- fmt.Errorf("listener %v is not created", name)
		return
	}

	if value.CancelFunc != nil {
		t.Result <- fmt.Errorf("listener %v is running", name)
		return
	}

	value.Listener = nil
	t.Result <- nil
}

func (s *ListenerController) handleStart(t *TaskStart) {
	defer close(t.Result)
	name := t.Name

	value, ok := s.state.Listeners.Load(name)
	if !ok {
		t.Result <- fmt.Errorf("listener %v does not exist", name)
		return
	}

	if value.Listener == nil {
		t.Result <- fmt.Errorf("listener %v is not created", name)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	value.CancelFunc = cancel
	go func() {
		value.Listener.Listen(ctx)
	}()

	t.Result <- nil
}

func (s *ListenerController) handleStop(t *TaskStop) {
	defer close(t.Result)
	name := t.Name

	value, ok := s.state.Listeners.Load(name)
	if !ok {
		t.Result <- fmt.Errorf("listener %v does not exist", name)
		return
	}

	if value.Listener == nil {
		t.Result <- fmt.Errorf("listener %v is not created", name)
		return
	}

	if value.CancelFunc == nil {
		t.Result <- fmt.Errorf("listener %v is not running", name)
		return
	}

	value.CancelFunc()
	value.CancelFunc = nil
	t.Result <- nil
}
