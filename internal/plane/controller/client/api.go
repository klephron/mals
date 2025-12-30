package client

import (
	"context"
	"fmt"
	"mals/internal/client"
	"mals/internal/listener"
	"sync"
)

func (s *ClientController) Shutdown() error {
	s.state.clients.Range(func(key client.Client, value *ClientValue) bool {
		s.ClientShutdown(key)
		return true
	})

	s.state.statusRW.RLock()
	cancel := s.state.statusCancel
	s.state.statusRW.RUnlock()

	cancel()

	return nil
}

func (s *ClientController) ClientOwn(client client.Client, listener listener.Listener) error {
	name := client.Name()

	if _, ok := s.state.clients.Load(client); ok {
		return fmt.Errorf("client %v is owned", name)
	}

	if err := s.plane.Listener().ListenerClientAdd(listener.Name(), client); err != nil {
		return err
	}

	s.state.clients.Store(client, &ClientValue{
		rw:         sync.RWMutex{},
		listener:   listener.Name(),
		cancelFunc: nil,
	})

	return nil
}

func (s *ClientController) ClientServe(client client.Client) error {
	name := client.Name()

	value, ok := s.state.clients.Load(client)
	if !ok {
		return fmt.Errorf("client %v does not exist", name)
	}

	value.rw.Lock()

	if value.cancelFunc != nil {
		value.rw.Unlock()
		return fmt.Errorf("client %v is running", name)
	}

	ctx, cancel := context.WithCancel(context.Background())

	value.cancelFunc = cancel
	value.rw.Unlock()

	go func() {
		client.Serve(ctx)
		s.ClientShutdown(client)
	}()

	return nil
}

func (s *ClientController) clientShutdown(client client.Client, notify bool) error {
	name := client.Name()

	value, ok := s.state.clients.Load(client)
	if !ok {
		return fmt.Errorf("client %v does not exist", name)
	}

	value.rw.Lock()

	if value.cancelFunc == nil {
		value.rw.Unlock()
		return fmt.Errorf("client %v is not running", name)
	}

	cancel := value.cancelFunc
	value.cancelFunc = nil

	value.rw.Unlock()

	cancel()

	s.state.clients.Delete(client)

	if notify {
		if err := s.plane.Listener().ListenerClientRemove(value.listener, client); err != nil {
			return err
		}
	}

	return nil
}

func (s *ClientController) ClientShutdown(client client.Client) error {
	return s.clientShutdown(client, true)
}

func (s *ClientController) ClientShutdownSilent(client client.Client) error {
	return s.clientShutdown(client, false)
}
