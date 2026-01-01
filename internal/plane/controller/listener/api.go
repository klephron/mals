package listener

import (
	"context"
	"fmt"
	"mals/internal/listener/api_tcp"
	"mals/internal/listener/lsp_tcp"
	"mals/internal/plane/controller"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

func statusErrorEq(name string, actual controller.ListenerStatus, expected controller.ListenerStatus) error {
	return fmt.Errorf("listener %v expected eq %v, got %v", name, expected, actual)
}

func statusErrorFlag(name string, actual controller.ListenerStatus, expected controller.ListenerStatus) error {
	return fmt.Errorf("listener %v expected flag %v, got %v", name, expected, actual)
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
		s.ListenerStop(key)
		s.ListenerDelete(key)
		return true
	})

	s.state.statusRW.RLock()
	cancel := s.state.statusCancel
	s.state.statusRW.RUnlock()

	cancel()

	return nil
}

func (s *ListenerController) ListenerStatus(name string) controller.ListenerStatus {
	value, _ := s.state.listeners.Load(name)
	return s.statusRW(value)
}

func (s *ListenerController) ListenerRegister(name string, cfg config.Listener) error {
	value, _ := s.state.listeners.Load(name)

	if status := s.statusRW(value); status != controller.ListenerAbsent {
		return statusErrorEq(name, status, controller.ListenerAbsent)
	}

	s.state.listeners.Store(name, &ListenerValue{
		rw:         sync.RWMutex{},
		config:     cfg,
		listener:   nil,
		cancelFunc: nil,
		clients:    xsync.NewMap[string, struct{}](),
	})

	return nil
}

func (s *ListenerController) ListenerUnregister(name string) error {
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

func (s *ListenerController) ListenerCreate(name string) error {
	value, _ := s.state.listeners.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status != controller.ListenerRegistered {
		return statusErrorEq(name, status, controller.ListenerRegistered)
	}

	switch value.config.Kind.(type) {
	case *config.ListenerKindApi:
		switch ipc := value.config.Ipc.(type) {

		case *config.ListenerIpcTcp:
			listener, err := api_tcp.NewListener(name, ipc.Port, s.plane)
			if err != nil {
				return err
			}
			value.listener = listener
			return nil
		}

	case *config.ListenerKindLsp:
		switch ipc := value.config.Ipc.(type) {

		case *config.ListenerIpcTcp:
			listener, err := lsp_tcp.NewListener(name, ipc.Port, s.plane)
			if err != nil {
				return err
			}
			value.listener = listener
			return nil
		}
	}

	return fmt.Errorf("unhandled listener %v kind=%v ipc=%v", value.config, value.config.Kind.Kind(), value.config.Ipc.Ipc())
}

func (s *ListenerController) ListenerDelete(name string) error {
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

func (s *ListenerController) ListenerStart(name string) error {
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
		listener.Run(ctx)
		s.ListenerStop(listener.Name())
	}()

	return nil
}

func (s *ListenerController) ListenerStop(name string) error {
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

	value.clients.Range(func(key string, value struct{}) bool {
		s.plane.ClientShutdown(key)
		return true
	})

	cancel := value.cancelFunc
	value.cancelFunc = nil
	value.rw.Unlock()

	cancel()

	return nil
}

func (s *ListenerController) ListenerClientAdd(name string, client string) error {
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
		return fmt.Errorf("listener %v client %v exists", name, client)
	}

	value.clients.Store(client, struct{}{})

	return nil
}

func (s *ListenerController) ListenerClientRemove(name string, client string) error {
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
		return fmt.Errorf("listener %v client %v does not exist", name, client)
	}

	return nil
}

func (s *ListenerController) ListenerGetConfig(name string) (config.Listener, error) {
	value, _ := s.state.listeners.Load(name)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.ListenerRegistered == 0 {
		return config.Listener{}, statusErrorFlag(name, status, controller.ListenerRegistered)
	}

	return value.config, nil
}
