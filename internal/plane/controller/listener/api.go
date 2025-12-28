package listener

import (
	"context"
	"fmt"
	"mals/internal/client"
	"mals/internal/listener/api_tcp"
	"mals/internal/listener/lsp_tcp"
	"mals/internal/plane/controller"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

func statusErrorEq(name string, actual controller.ListenerStatus, expected controller.ListenerStatus) error {
	return fmt.Errorf("Listener %v expected eq %v, got %v", name, expected, actual)
}

func statusErrorFlag(name string, actual controller.ListenerStatus, expected controller.ListenerStatus) error {
	return fmt.Errorf("Listener %v expected flag %v, got %v", name, expected, actual)
}

func (s *ListenerController) status(value *ListenerValue) controller.ListenerStatus {
	status := controller.ListenerAbsent

	if value != nil {
		status |= controller.ListenerRegistered

		if value.listener != nil {
			status |= controller.ListenerCreated
		}
		if value.cancelFunc != nil {
			status |= controller.ListenerStarted
		}
	}

	return status
}

func (s *ListenerController) statusRW(value *ListenerValue) controller.ListenerStatus {
	status := controller.ListenerAbsent

	if value != nil {
		status |= controller.ListenerRegistered

		value.rw.RLock()

		if value.listener != nil {
			status |= controller.ListenerCreated
		}
		if value.cancelFunc != nil {
			status |= controller.ListenerStarted
		}

		value.rw.RUnlock()
	}

	return status
}

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

func (s *ListenerController) Status(name string) controller.ListenerStatus {
	value, _ := s.state.listeners.Load(name)
	return s.statusRW(value)
}

func (s *ListenerController) Register(cfg config.Listener) error {
	name := cfg.Name()

	value, _ := s.state.listeners.Load(name)

	if status := s.statusRW(value); status != controller.ListenerAbsent {
		return statusErrorEq(name, status, controller.ListenerAbsent)
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
	value, _ := s.state.listeners.Load(name)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status != controller.ListenerRegistered {
		return statusErrorEq(name, status, controller.ListenerRegistered)
	}

	s.state.listeners.Delete(name)

	return nil
}

func (s *ListenerController) Create(name string) error {
	value, _ := s.state.listeners.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status != controller.ListenerRegistered {
		return statusErrorEq(name, status, controller.ListenerRegistered)
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
	value, _ := s.state.listeners.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status != controller.ListenerRegistered|controller.ListenerCreated {
		return statusErrorEq(name, status, controller.ListenerRegistered|controller.ListenerCreated)
	}

	value.listener = nil

	return nil
}

func (s *ListenerController) Start(name string) error {
	value, _ := s.state.listeners.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status != controller.ListenerRegistered|controller.ListenerCreated {
		return statusErrorEq(name, status, controller.ListenerRegistered|controller.ListenerCreated)
	}

	ctx, cancel := context.WithCancel(context.Background())

	value.cancelFunc = cancel
	listener := value.listener

	go func() {
		listener.Listen(ctx)
		s.Stop(listener.Name())
	}()

	return nil
}

func (s *ListenerController) Stop(name string) error {
	value, _ := s.state.listeners.Load(name)

	if value != nil {
		value.rw.Lock()
	}

	if status := s.status(value); status&controller.ListenerStarted == 0 {
		if value != nil {
			value.rw.Unlock()
		}
		return statusErrorFlag(name, status, controller.ListenerStarted)
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
	value, _ := s.state.listeners.Load(name)

	if value != nil {
		value.rw.RLock()
	}

	if status := s.status(value); status&controller.ListenerStarted == 0 {
		if value != nil {
			value.rw.RUnlock()
		}
		return statusErrorFlag(name, status, controller.ListenerStarted)
	}

	value.rw.RUnlock()

	_, ok := value.clients.Load(client)
	if ok {
		return fmt.Errorf("listener %v client %v exists", name, client.Name())
	}

	value.clients.Store(client, struct{}{})

	return nil
}

func (s *ListenerController) ClientRemove(name string, client client.Client) error {
	value, _ := s.state.listeners.Load(name)

	if value != nil {
		value.rw.RLock()
	}

	if status := s.status(value); status&controller.ListenerStarted == 0 {
		if value != nil {
			value.rw.RUnlock()
		}
		return statusErrorFlag(name, status, controller.ListenerStarted)
	}

	value.rw.RUnlock()

	_, ok := value.clients.LoadAndDelete(client)
	if !ok {
		return fmt.Errorf("listener %v client %v does not exist", name, client.Name())
	}

	return nil
}
