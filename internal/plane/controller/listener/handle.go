package listener

import (
	"context"
	"fmt"
	"mals/internal/client"
	"mals/internal/listener"
	"mals/internal/listener/api/tcp"
	"mals/internal/listener/lsp/tcp"
	"mals/pkg/config"

	"github.com/puzpuzpuz/xsync/v4"
)

func (s *ListenerController) handleShutdown(t *TaskShutdown) {
	defer close(t.result)

	s.state.listeners.Range(func(key string, value *ListenerValue) bool {
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

func (s *ListenerController) handleRegister(t *TaskRegister) {
	defer close(t.result)
	name := t.config.Name()

	if _, ok := s.state.listeners.Load(name); ok {
		t.result <- fmt.Errorf("listener %v exists", name)
		return
	}

	switch config := t.config.(type) {
	case *config.ListenerTcp:
		s.state.listeners.Store(name, &ListenerValue{
			config:     config,
			listener:   nil,
			cancelFunc: nil,
			clients:    xsync.NewMap[client.Client, struct{}](),
		})
		t.result <- nil

	default:
		t.result <- fmt.Errorf("unhandled listener %T %v", config, config)
	}
}

func (s *ListenerController) handleUnregister(t *TaskUnregister) {
	defer close(t.result)
	name := t.name

	value, ok := s.state.listeners.Load(name)
	if !ok {
		t.result <- fmt.Errorf("listener %v does not exist", name)
		return
	}
	if value.listener != nil {
		t.result <- fmt.Errorf("listener %v is already created", name)
		return
	}

	s.state.listeners.Delete(name)

	t.result <- nil
}

func (s *ListenerController) handleCreate(t *TaskCreate) {
	defer close(t.result)
	name := t.name

	value, ok := s.state.listeners.Load(name)
	if !ok {
		t.result <- fmt.Errorf("listener %v does not exist", name)
		return
	}
	if value.listener != nil {
		t.result <- fmt.Errorf("listener %v is already created", name)
		return
	}

	switch config := value.config.(type) {
	case *config.ListenerTcp:
		switch config.Kind() {

		case apitcp.Kind():
			listener, err := apitcp.NewListener(name, config.Port, s.plane)
			if err != nil {
				t.result <- err
				return
			}
			value.listener = listener
			t.result <- nil

		case lsptcp.Kind():
			listener, err := lsptcp.NewListener(name, config.Port, s.plane)
			if err != nil {
				t.result <- err
				return
			}
			value.listener = listener
			t.result <- nil

		default:
			t.result <- fmt.Errorf("unhandled listener %T %v", config, config)
		}

	default:
		t.result <- fmt.Errorf("unhandled listener %T %v", config, config)
	}
}

func (s *ListenerController) handleDelete(t *TaskDelete) {
	defer close(t.result)
	name := t.name

	value, ok := s.state.listeners.Load(name)
	if !ok {
		t.result <- fmt.Errorf("listener %v does not exist", name)
		return
	}
	if value.listener == nil {
		t.result <- fmt.Errorf("listener %v is not created", name)
		return
	}
	if value.cancelFunc != nil {
		t.result <- fmt.Errorf("listener %v is running", name)
		return
	}

	value.listener = nil
	t.result <- nil
}

func (s *ListenerController) handleStart(t *TaskStart) {
	defer close(t.result)
	name := t.name

	value, ok := s.state.listeners.Load(name)
	if !ok {
		t.result <- fmt.Errorf("listener %v does not exist", name)
		return
	}
	if value.listener == nil {
		t.result <- fmt.Errorf("listener %v is not created", name)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	value.cancelFunc = cancel

	go func(listener listener.Listener) {
		listener.Listen(ctx)
		// avoid race conditions
		s.Stop(listener.Name())
	}(value.listener)

	t.result <- nil
}

func (s *ListenerController) handleStop(t *TaskStop) {
	defer close(t.result)
	name := t.name

	value, ok := s.state.listeners.Load(name)
	if !ok {
		t.result <- fmt.Errorf("listener %v does not exist", name)
		return
	}
	if value.listener == nil {
		t.result <- fmt.Errorf("listener %v is not created", name)
		return
	}
	if value.cancelFunc == nil {
		t.result <- fmt.Errorf("listener %v is not running", name)
		return
	}

	value.clients.Range(func(key client.Client, value struct{}) bool {
		s.plane.Client().Stop(key)
		s.plane.Client().DeleteSilent(key)
		return true
	})

	value.cancelFunc()
	value.cancelFunc = nil

	t.result <- nil
}

func (s *ListenerController) handleClientAdd(t *TaskClientAdd) {
	defer close(t.result)
	name := t.name

	value, ok := s.state.listeners.Load(name)
	if !ok {
		t.result <- fmt.Errorf("listener %v does not exist", name)
		return
	}
	if value.listener == nil {
		t.result <- fmt.Errorf("listener %v is not created", name)
		return
	}
	if value.cancelFunc == nil {
		t.result <- fmt.Errorf("listener %v is not running", name)
		return
	}

	_, ok = value.clients.Load(t.client)
	if ok {
		t.result <- fmt.Errorf("listener %v client %v exists", name, t.client.Name())
		return
	}

	value.clients.Store(t.client, struct{}{})

	t.result <- nil
}

func (s *ListenerController) handleClientRemove(t *TaskClientRemove) {
	defer close(t.result)
	name := t.name

	value, ok := s.state.listeners.Load(name)
	if !ok {
		t.result <- fmt.Errorf("listener %v does not exist", name)
		return
	}
	if value.listener == nil {
		t.result <- fmt.Errorf("listener %v is not created", name)
		return
	}
	if value.cancelFunc == nil {
		t.result <- fmt.Errorf("listener %v is not running", name)
		return
	}

	_, ok = value.clients.LoadAndDelete(t.client)
	if !ok {
		t.result <- fmt.Errorf("listener %v client %v does not exist", name, t.client.Name())
		return
	}

	t.result <- nil
}
