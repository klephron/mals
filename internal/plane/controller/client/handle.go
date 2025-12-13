package client

import (
	"context"
	"fmt"
	"mals/internal/client"
	"mals/internal/plane/state"
)

func (s *ClientController) handleShutdown(t *TaskShutdown) {
	defer close(t.Result)

	s.state.Clients.Range(func(key client.Client, value *state.ClientValue) bool {
		ts := &TaskStop{TaskGeneric: NewTask(), Client: key}
		s.handleStop(ts)
		<-ts.Result

		td := &TaskDelete{TaskGeneric: NewTask(), Client: key}
		s.handleDelete(td)
		<-td.Result

		return true
	})

	t.Result <- nil
}

func (s *ClientController) handleOwn(t *TaskOwn) {
	defer close(t.Result)
	name := t.Client.Name()

	if _, ok := s.state.Clients.Load(t.Client); ok {
		t.Result <- fmt.Errorf("client %v is owned", name)
		return
	}

	if err := s.plane.Listener().ClientAdd(t.Listener, t.Client); err != nil {
		t.Result <- err
		return
	}

	s.state.Clients.Store(t.Client, &state.ClientValue{
		Listener:   t.Listener,
		CancelFunc: nil,
	})

	t.Result <- nil
}

func (s *ClientController) handleDelete(t *TaskDelete) {
	defer close(t.Result)
	name := t.Client.Name()

	value, ok := s.state.Clients.Load(t.Client)
	if !ok {
		t.Result <- fmt.Errorf("client %v does not exist", name)
		return
	}
	if value.CancelFunc != nil {
		t.Result <- fmt.Errorf("client %v is running", name)
		return
	}

	if t.Notify {
		if err := s.plane.Listener().ClientRemove(value.Listener, t.Client); err != nil {
			t.Result <- err
			return
		}
	}

	s.state.Clients.Delete(t.Client)

	t.Result <- nil
}

func (s *ClientController) handleStart(t *TaskStart) {
	defer close(t.Result)
	name := t.Client.Name()

	value, ok := s.state.Clients.Load(t.Client)
	if !ok {
		t.Result <- fmt.Errorf("client %v does not exist", name)
		return
	}
	if value.CancelFunc != nil {
		t.Result <- fmt.Errorf("client %v is running", name)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	value.CancelFunc = cancel
	go func(client client.Client) {
		client.Serve(ctx)
		// NOTE: to avoid race conditions
		s.Stop(client)
		// NOTE:  automatically delete when stopped
		s.Delete(client)
	}(t.Client)

	t.Result <- nil
}

func (s *ClientController) handleStop(t *TaskStop) {
	defer close(t.Result)
	name := t.Client.Name()

	value, ok := s.state.Clients.Load(t.Client)
	if !ok {
		t.Result <- fmt.Errorf("client %v does not exist", name)
		return
	}
	if value.CancelFunc == nil {
		t.Result <- fmt.Errorf("client %v is not running", name)
		return
	}

	value.CancelFunc()
	value.CancelFunc = nil

	t.Result <- nil
}
