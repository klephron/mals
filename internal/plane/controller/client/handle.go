package client

import (
	"context"
	"fmt"
	"mals/internal/client"
)

func (s *ClientController) handleShutdown(t *TaskShutdown) {
	defer close(t.result)

	s.state.clients.Range(func(key client.Client, value *ClientValue) bool {
		ts := &TaskStop{TaskGeneric: newTask(), client: key}
		s.handleStop(ts)
		<-ts.result

		td := &TaskDelete{TaskGeneric: newTask(), client: key}
		s.handleDelete(td)
		<-td.result

		return true
	})

	t.result <- nil
}

func (s *ClientController) handleOwn(t *TaskOwn) {
	defer close(t.result)
	name := t.client.Name()

	if _, ok := s.state.clients.Load(t.client); ok {
		t.result <- fmt.Errorf("client %v is owned", name)
		return
	}

	if err := s.plane.Listener().ClientAdd(t.listener, t.client); err != nil {
		t.result <- err
		return
	}

	s.state.clients.Store(t.client, &ClientValue{
		listener:   t.listener,
		cancelFunc: nil,
	})

	t.result <- nil
}

func (s *ClientController) handleDelete(t *TaskDelete) {
	defer close(t.result)
	name := t.client.Name()

	value, ok := s.state.clients.Load(t.client)
	if !ok {
		t.result <- fmt.Errorf("client %v does not exist", name)
		return
	}
	if value.cancelFunc != nil {
		t.result <- fmt.Errorf("client %v is running", name)
		return
	}

	if t.notify {
		if err := s.plane.Listener().ClientRemove(value.listener, t.client); err != nil {
			t.result <- err
			return
		}
	}

	s.state.clients.Delete(t.client)

	t.result <- nil
}

func (s *ClientController) handleStart(t *TaskStart) {
	defer close(t.result)
	name := t.client.Name()

	value, ok := s.state.clients.Load(t.client)
	if !ok {
		t.result <- fmt.Errorf("client %v does not exist", name)
		return
	}
	if value.cancelFunc != nil {
		t.result <- fmt.Errorf("client %v is running", name)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	value.cancelFunc = cancel
	go func(client client.Client) {
		client.Serve(ctx)
		// NOTE: to avoid race conditions
		s.Stop(client)
		// NOTE:  automatically delete when stopped
		s.Delete(client)
	}(t.client)

	t.result <- nil
}

func (s *ClientController) handleStop(t *TaskStop) {
	defer close(t.result)
	name := t.client.Name()

	value, ok := s.state.clients.Load(t.client)
	if !ok {
		t.result <- fmt.Errorf("client %v does not exist", name)
		return
	}
	if value.cancelFunc == nil {
		t.result <- fmt.Errorf("client %v is not running", name)
		return
	}

	value.cancelFunc()
	value.cancelFunc = nil

	t.result <- nil
}
