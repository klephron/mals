package listener

import (
	"context"
	"fmt"
	"mals/internal/client"
	"mals/internal/listener/api/tcp"
	"mals/internal/listener/lsp/tcp"
	"mals/internal/plane/state"
	"mals/pkg/config"

	"github.com/puzpuzpuz/xsync/v4"
)

func (s *ListenerController) handleShutdown(t *TaskShutdown) error {
	defer close(t.Result)
	s.state.Listeners.Range(func(key string, value *state.ListenerValue) bool {
		ts := &TaskStop{TaskGeneric: NewTaskSingle(), Name: key}
		s.handleStop(ts)
		<-ts.Result

		td := &TaskDelete{TaskGeneric: NewTaskSingle(), Name: key}
		s.handleDelete(td)
		<-td.Result

		return true
	})
	return nil
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
			Clients:    xsync.NewMap[client.Client, struct{}](),
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
		case apitcp.Kind():
			if listener, err := apitcp.NewListener(name, config.Port, s.plane); err != nil {
				t.Result <- err
			} else {
				value.Listener = listener
				t.Result <- nil
			}

		case lsptcp.Kind():
			if listener, err := lsptcp.NewListener(name, config.Port, s.plane); err != nil {
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

	value.Clients.Range(func(key client.Client, value struct{}) bool {
		if err := s.plane.Client().Stop(key); err != nil {
			s.plane.Log().Errorf("%v", err)
		}
		if err := s.plane.Client().Delete(key); err != nil {
			s.plane.Log().Errorf("%v", err)
		}
		return true
	})

	value.CancelFunc()
	value.CancelFunc = nil

	t.Result <- nil
}

func (s *ListenerController) handleClientAdd(t *TaskClientAdd) {
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

	_, ok = value.Clients.Load(t.Client)
	if ok {
		t.Result <- fmt.Errorf("listener %v client %v exists", name, t.Client.Name())
	}

	value.Clients.Store(t.Client, struct{}{})
}

func (s *ListenerController) handleClientRemove(t *TaskClientRemove) {
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

	_, ok = value.Clients.LoadAndDelete(t.Client)
	if !ok {
		t.Result <- fmt.Errorf("listener %v client %v does not exist", name, t.Client.Name())
	}
}
