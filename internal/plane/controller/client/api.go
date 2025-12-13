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
		s.Stop(key)
		s.Delete(key)
		return true
	})

	s.state.statusRW.RLock()
	cancel := s.state.statusCancel
	s.state.statusRW.RUnlock()

	cancel()

	return nil
}

func (s *ClientController) Own(client client.Client, listener listener.Listener) error {
	name := client.Name()

	if _, ok := s.state.clients.Load(client); ok {
		return fmt.Errorf("client %v is owned", name)
	}

	if err := s.plane.Listener().ClientAdd(listener.Name(), client); err != nil {
		return err
	}

	s.state.clients.Store(client, &ClientValue{
		rw:         sync.RWMutex{},
		listener:   listener.Name(),
		cancelFunc: nil,
	})

	return nil
}

func (s *ClientController) delete(client client.Client, notify bool) error {
	name := client.Name()

	value, ok := s.state.clients.Load(client)
	if !ok {
		return fmt.Errorf("client %v does not exist", name)
	}
	if value.cancelFunc != nil {
		return fmt.Errorf("client %v is running", name)
	}

	s.state.clients.Delete(client)

	if notify {
		if err := s.plane.Listener().ClientRemove(value.listener, client); err != nil {
			return err
		}
	}

	return nil
}

func (s *ClientController) Delete(client client.Client) error {
	return s.delete(client, true)
}

func (s *ClientController) DeleteSilent(client client.Client) error {
	return s.delete(client, false)
}

func (s *ClientController) Start(client client.Client) error {
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
		s.Stop(client)
		s.Delete(client)
	}()

	return nil
}

func (s *ClientController) Stop(client client.Client) error {
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

	return nil
}
