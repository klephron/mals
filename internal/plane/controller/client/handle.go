package client

import (
	"context"
	"fmt"
	"mals/internal/plane/state"
)

func (s *ClientController) handleShutdown(_ *TaskShutdown) error {
	return nil
}

func (s *ClientController) handleOwn(t *TaskOwn) {
	defer close(t.Result)
	name := t.Client.Name()

	if _, ok := s.state.Clients.Load(t.Client); ok {
		t.Result <- fmt.Errorf("client %v is owned", name)
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
	go func() {
		t.Client.Serve(ctx)
	}()
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
