package client

import (
	"context"
	"fmt"
	"mals/internal/client"
	"mals/internal/listener"
	"mals/internal/plane/controller"
	"sync"
)

func statusErrorEq(client string, actual controller.ClientStatus, expected controller.ClientStatus) error {
	return fmt.Errorf("client %v expected eq %v, got %v", client, expected, actual)
}

func statusErrorFlag(client string, actual controller.ClientStatus, expected controller.ClientStatus) error {
	return fmt.Errorf("client %v expected flag %v, got %v", client, expected, actual)
}

func (s *ClientController) status(value *ClientValue) controller.ClientStatus {
	status := controller.ClientAbsent

	if value != nil {
		status |= controller.ClientCreated

		if value.cancelFunc != nil {
			status |= controller.ClientStarted
		}
	}

	return status
}

func (s *ClientController) statusRW(value *ClientValue) controller.ClientStatus {
	status := controller.ClientAbsent

	if value != nil {
		status |= controller.ClientCreated

		if value.cancelFunc != nil {
			status |= controller.ClientStarted
		}
	}

	return status
}

func (s *ClientController) Shutdown() error {
	s.state.clients.Range(func(key string, value *ClientValue) bool {
		s.ClientShutdown(key)
		return true
	})

	s.state.statusRW.RLock()
	cancel := s.state.statusCancel
	s.state.statusRW.RUnlock()

	cancel()

	return nil
}

func (s *ClientController) ClientStatus(name string) controller.ClientStatus {
	value, _ := s.state.clients.Load(name)
	return s.statusRW(value)
}

func (s *ClientController) ClientOwn(name string, client client.Client, listener listener.Listener) error {
	value, _ := s.state.clients.Load(name)

	if status := s.statusRW(value); status != controller.ClientAbsent {
		return statusErrorEq(name, status, controller.ClientAbsent)
	}

	if err := s.plane.ListenerClientAdd(listener.Name(), name); err != nil {
		return err
	}

	s.state.clients.Store(name, &ClientValue{
		rw:         sync.RWMutex{},
		client:     client,
		listener:   listener.Name(),
		cancelFunc: nil,
	})

	return nil
}

func (s *ClientController) ClientServe(name string) error {
	value, _ := s.state.clients.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status != controller.ClientCreated {
		return statusErrorEq(name, status, controller.ClientCreated)
	}

	ctx, cancel := context.WithCancel(context.Background())

	value.cancelFunc = cancel

	go func() {
		value.client.Serve(ctx)
		s.ClientShutdown(name)
	}()

	return nil
}

func (s *ClientController) clientShutdown(name string, notify bool) error {
	value, _ := s.state.clients.Load(name)

	if value != nil {
		value.rw.Lock()
	}

	if status := s.status(value); status&controller.ClientStarted == 0 {
		return statusErrorFlag(name, status, controller.ClientStarted)
	}

	cancel := value.cancelFunc
	value.cancelFunc = nil

	value.rw.Unlock()

	cancel()

	s.state.clients.Delete(name)

	if notify {
		if err := s.plane.ListenerClientRemove(value.listener, name); err != nil {
			return err
		}
	}

	return nil
}

func (s *ClientController) ClientShutdown(name string) error {
	return s.clientShutdown(name, true)
}

func (s *ClientController) ClientShutdownSilent(name string) error {
	return s.clientShutdown(name, false)
}

func (s *ClientController) ClientGetListener(name string) (string, error) {
	value, _ := s.state.clients.Load(name)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.ClientCreated == 0 {
		return "", statusErrorFlag(name, status, controller.ClientCreated)
	}

	return value.listener, nil
}
