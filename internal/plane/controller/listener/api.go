package listener

import (
	"context"
	"fmt"
	"mals/internal/client"
	"mals/internal/listener/api_tcp"
	"mals/internal/listener/lsp_tcp"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

func (s *ListenerController) Shutdown() error {
	s.state.listeners.Range(func(key string, value *ListenerValue) bool {
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

func (s *ListenerController) Register(cfg config.Listener) error {
	name := cfg.Name()

	if _, ok := s.state.listeners.Load(name); ok {
		return fmt.Errorf("listener %v exists", name)
	}

	switch config := cfg.(type) {
	case *config.ListenerTcp:
		s.state.listeners.Store(name, &ListenerValue{
			rw:         sync.RWMutex{},
			config:     config,
			listener:   nil,
			cancelFunc: nil,
			clients:    xsync.NewMap[client.Client, struct{}](),
		})

	default:
		return fmt.Errorf("unhandled listener %T %v", config, config)
	}

	return nil
}

func (s *ListenerController) Unregister(name string) error {
	value, ok := s.state.listeners.Load(name)
	if !ok {
		return fmt.Errorf("listener %v does not exist", name)
	}

	value.rw.RLock()
	if value.listener != nil {
		value.rw.RUnlock()
		return fmt.Errorf("listener %v is already created", name)
	}

	s.state.listeners.Delete(name)
	value.rw.RUnlock()

	return nil
}

func (s *ListenerController) Create(name string) error {
	value, ok := s.state.listeners.Load(name)
	if !ok {
		return fmt.Errorf("listener %v does not exist", name)
	}

	value.rw.Lock()
	defer value.rw.Unlock()

	if value.listener != nil {
		return fmt.Errorf("listener %v is already created", name)
	}

	switch config := value.config.(type) {
	case *config.ListenerTcp:
		switch config.Kind() {

		case api_tcp.Kind():
			listener, err := api_tcp.NewListener(name, config.Port, s.plane)
			if err != nil {
				return err
			}
			value.listener = listener

		case lsp_tcp.Kind():
			listener, err := lsp_tcp.NewListener(name, config.Port, s.plane)
			if err != nil {
				return err
			}
			value.listener = listener

		default:
			return fmt.Errorf("unhandled listener %T %v", config, config)
		}

	default:
		return fmt.Errorf("unhandled listener %T %v", config, config)
	}

	return nil
}

func (s *ListenerController) Delete(name string) error {
	value, ok := s.state.listeners.Load(name)
	if !ok {
		return fmt.Errorf("listener %v does not exist", name)
	}

	value.rw.Lock()

	if value.listener == nil {
		value.rw.Unlock()
		return fmt.Errorf("listener %v is not created", name)
	}
	if value.cancelFunc != nil {
		value.rw.Unlock()
		return fmt.Errorf("listener %v is running", name)
	}

	value.listener = nil
	value.rw.Unlock()

	return nil
}

func (s *ListenerController) Start(name string) error {
	value, ok := s.state.listeners.Load(name)
	if !ok {
		return fmt.Errorf("listener %v does not exist", name)
	}

	value.rw.Lock()

	if value.listener == nil {
		value.rw.Unlock()
		return fmt.Errorf("listener %v is not created", name)
	}

	ctx, cancel := context.WithCancel(context.Background())

	value.cancelFunc = cancel
	listener := value.listener
	value.rw.Unlock()

	go func() {
		listener.Listen(ctx)
		s.Stop(listener.Name())
	}()

	return nil
}

func (s *ListenerController) Stop(name string) error {
	value, ok := s.state.listeners.Load(name)
	if !ok {
		return fmt.Errorf("listener %v does not exist", name)
	}

	value.rw.Lock()

	if value.listener == nil {
		value.rw.Unlock()
		return fmt.Errorf("listener %v is not created", name)
	}
	if value.cancelFunc == nil {
		value.rw.Unlock()
		return fmt.Errorf("listener %v is not running", name)
	}

	value.clients.Range(func(key client.Client, value struct{}) bool {
		s.plane.Client().Stop(key)
		s.plane.Client().DeleteSilent(key)
		return true
	})

	cancel := value.cancelFunc
	value.cancelFunc = nil
	value.rw.Unlock()

	cancel()

	return nil
}

func (s *ListenerController) ClientAdd(name string, client client.Client) error {
	value, ok := s.state.listeners.Load(name)
	if !ok {
		return fmt.Errorf("listener %v does not exist", name)
	}

	value.rw.RLock()

	if value.listener == nil {
		value.rw.RUnlock()
		return fmt.Errorf("listener %v is not created", name)
	}
	if value.cancelFunc == nil {
		value.rw.RUnlock()
		return fmt.Errorf("listener %v is not running", name)
	}

	value.rw.RUnlock()

	_, ok = value.clients.Load(client)
	if ok {
		return fmt.Errorf("listener %v client %v exists", name, client.Name())
	}

	value.clients.Store(client, struct{}{})

	return nil
}

func (s *ListenerController) ClientRemove(name string, client client.Client) error {
	value, ok := s.state.listeners.Load(name)
	if !ok {
		return fmt.Errorf("listener %v does not exist", name)
	}

	value.rw.RLock()

	if value.listener == nil {
		value.rw.RUnlock()
		return fmt.Errorf("listener %v is not created", name)
	}
	if value.cancelFunc == nil {
		value.rw.RUnlock()
		return fmt.Errorf("listener %v is not running", name)
	}

	value.rw.RUnlock()

	_, ok = value.clients.LoadAndDelete(client)
	if !ok {
		return fmt.Errorf("listener %v client %v does not exist", name, client.Name())
	}

	return nil
}
