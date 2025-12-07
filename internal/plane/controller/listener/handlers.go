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
		s.handleStop(&EventStop{Name: key})
		s.handleDelete(&EventDelete{Name: key})
		return true
	})
}

func (s *ListenerController) handleTerminate(_ *event.EventTerminate) {
	s.state.Listeners.Range(func(key string, value *state.ListenerValue) bool {
		s.handleDelete(&EventDelete{Name: key})
		return true
	})
}

func (s *ListenerController) handleRegister(e *EventRegister) {
	name := e.Config.Name()

	if _, ok := s.state.Listeners.Load(name); ok {
		e.Error = fmt.Errorf("listener %v exists", name)
		return
	}

	switch config := e.Config.(type) {
	case *config.ListenerTcp:
		s.state.Listeners.Store(name, &state.ListenerValue{
			Config:     config,
			Listener:   nil,
			CancelFunc: nil,
		})
		e.Error = nil

	default:
		e.Error = fmt.Errorf("unhandled listener %T %v", config, config)
	}
}

func (s *ListenerController) handleUnregister(e *EventUnregister) {
	name := e.Name

	value, ok := s.state.Listeners.Load(name)
	if !ok {
		e.Error = fmt.Errorf("listener %v does not exist", name)
		return
	}

	if value.Listener != nil {
		e.Error = fmt.Errorf("listener %v is already created", name)
		return
	}

	s.state.Listeners.Delete(name)
	e.Error = nil
}

func (s *ListenerController) handleCreate(e *EventCreate) {
	name := e.Name

	value, ok := s.state.Listeners.Load(name)
	if !ok {
		e.Error = fmt.Errorf("listener %v does not exist", name)
		return
	}

	if value.Listener != nil {
		e.Error = fmt.Errorf("listener %v is already created", name)
		return
	}

	switch config := value.Config.(type) {
	case *config.ListenerTcp:
		switch config.Kind() {
		case lsp.Kind():
			if listener, err := lsp.New(config.Port, s.plane); err != nil {
				e.Error = err
			} else {
				value.Listener = listener
				e.Error = nil
			}

		case api.Kind():
			if listener, err := api.New(config.Port, s.plane); err != nil {
				e.Error = err
			} else {
				value.Listener = listener
				e.Error = nil
			}

		default:
			e.Error = fmt.Errorf("unhandled listener %T %v", config, config)
		}

	default:
		e.Error = fmt.Errorf("unhandled listener %T %v", config, config)
	}
}

func (s *ListenerController) handleDelete(e *EventDelete) {
	name := e.Name

	value, ok := s.state.Listeners.Load(name)
	if !ok {
		e.Error = fmt.Errorf("listener %v does not exist", name)
		return
	}

	if value.Listener == nil {
		e.Error = fmt.Errorf("listener %v is not created", name)
		return
	}

	if value.CancelFunc != nil {
		e.Error = fmt.Errorf("listener %v is running", name)
		return
	}

	value.Listener = nil
	e.Error = nil
}

func (s *ListenerController) handleStart(e *EventStart) {
	name := e.Name

	value, ok := s.state.Listeners.Load(name)
	if !ok {
		e.Error = fmt.Errorf("listener %v does not exist", name)
		return
	}

	if value.Listener == nil {
		e.Error = fmt.Errorf("listener %v is not created", name)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	value.Listener.Listen(ctx)
	value.CancelFunc = cancel
	e.Error = nil
}

func (s *ListenerController) handleStop(e *EventStop) {
	name := e.Name

	value, ok := s.state.Listeners.Load(name)
	if !ok {
		e.Error = fmt.Errorf("listener %v does not exist", name)
		return
	}

	if value.Listener == nil {
		e.Error = fmt.Errorf("listener %v is not created", name)
		return
	}

	if value.CancelFunc == nil {
		e.Error = fmt.Errorf("listener %v is not running", name)
		return
	}

	value.CancelFunc()
	value.CancelFunc = nil
	e.Error = nil
}
