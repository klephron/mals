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
		s.handleStop(&EventStop{EventGeneric: NewEventSingle(), Name: key})
		s.handleDelete(&EventDelete{EventGeneric: NewEventSingle(), Name: key})
		return true
	})
}

func (s *ListenerController) handleTerminate(_ *event.EventTerminate) {
	s.state.Listeners.Range(func(key string, value *state.ListenerValue) bool {
		s.state.Listeners.Delete(key)
		return true
	})
}

func (s *ListenerController) handleRegister(e *EventRegister) {
	defer close(e.Result)
	name := e.Config.Name()

	if _, ok := s.state.Listeners.Load(name); ok {
		e.Result <- fmt.Errorf("listener %v exists", name)
		return
	}

	switch config := e.Config.(type) {
	case *config.ListenerTcp:
		s.state.Listeners.Store(name, &state.ListenerValue{
			Config:     config,
			Listener:   nil,
			CancelFunc: nil,
		})
		e.Result <- nil

	default:
		e.Result <- fmt.Errorf("unhandled listener %T %v", config, config)
	}
}

func (s *ListenerController) handleUnregister(e *EventUnregister) {
	defer close(e.Result)
	name := e.Name

	value, ok := s.state.Listeners.Load(name)
	if !ok {
		e.Result <- fmt.Errorf("listener %v does not exist", name)
		return
	}

	if value.Listener != nil {
		e.Result <- fmt.Errorf("listener %v is already created", name)
		return
	}

	s.state.Listeners.Delete(name)
	e.Result <- nil
}

func (s *ListenerController) handleCreate(e *EventCreate) {
	defer close(e.Result)
	name := e.Name

	value, ok := s.state.Listeners.Load(name)
	if !ok {
		e.Result <- fmt.Errorf("listener %v does not exist", name)
		return
	}

	if value.Listener != nil {
		e.Result <- fmt.Errorf("listener %v is already created", name)
		return
	}

	switch config := value.Config.(type) {
	case *config.ListenerTcp:
		switch config.Kind() {
		case lsp.Kind():
			if listener, err := lsp.New(name, config.Port, s.plane); err != nil {
				e.Result <- err
			} else {
				value.Listener = listener
				e.Result <- nil
			}

		case api.Kind():
			if listener, err := api.New(name, config.Port, s.plane); err != nil {
				e.Result <- err
			} else {
				value.Listener = listener
				e.Result <- nil
			}

		default:
			e.Result <- fmt.Errorf("unhandled listener %T %v", config, config)
		}

	default:
		e.Result <- fmt.Errorf("unhandled listener %T %v", config, config)
	}
}

func (s *ListenerController) handleDelete(e *EventDelete) {
	defer close(e.Result)
	name := e.Name

	value, ok := s.state.Listeners.Load(name)
	if !ok {
		e.Result <- fmt.Errorf("listener %v does not exist", name)
		return
	}

	if value.Listener == nil {
		e.Result <- fmt.Errorf("listener %v is not created", name)
		return
	}

	if value.CancelFunc != nil {
		e.Result <- fmt.Errorf("listener %v is running", name)
		return
	}

	value.Listener = nil
	e.Result <- nil
}

func (s *ListenerController) handleStart(e *EventStart) {
	defer close(e.Result)
	name := e.Name

	value, ok := s.state.Listeners.Load(name)
	if !ok {
		e.Result <- fmt.Errorf("listener %v does not exist", name)
		return
	}

	if value.Listener == nil {
		e.Result <- fmt.Errorf("listener %v is not created", name)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	value.CancelFunc = cancel
	go func() {
		value.Listener.Listen(ctx)
	}()

	e.Result <- nil
}

func (s *ListenerController) handleStop(e *EventStop) {
	defer close(e.Result)
	name := e.Name

	value, ok := s.state.Listeners.Load(name)
	if !ok {
		e.Result <- fmt.Errorf("listener %v does not exist", name)
		return
	}

	if value.Listener == nil {
		e.Result <- fmt.Errorf("listener %v is not created", name)
		return
	}

	if value.CancelFunc == nil {
		e.Result <- fmt.Errorf("listener %v is not running", name)
		return
	}

	value.CancelFunc()
	value.CancelFunc = nil
	e.Result <- nil
}
