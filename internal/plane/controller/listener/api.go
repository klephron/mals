package listener

import (
	"context"
	"fmt"
	"mals/internal/listener"
	"mals/internal/listener/api"
	"mals/internal/listener/lsp"
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

func (s *ListenerController) status(value *Listener) controller.ListenerStatus {
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

func (s *ListenerController) statusRW(value *Listener) controller.ListenerStatus {
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

func clientStatusErrorEq(name string, clientName string, actual controller.ClientStatus, expected controller.ClientStatus) error {
	return fmt.Errorf("listener %v client %v expected eq %v, got %v", name, clientName, expected, actual)
}

func clientStatusErrorFlag(name string, clientName string, actual controller.ClientStatus, expected controller.ClientStatus) error {
	return fmt.Errorf("listener %v client %v expected flag %v, got %v", name, clientName, expected, actual)
}

func (s *ListenerController) lspClientStatus(value *ListenerLspClient) controller.ClientStatus {
	status := controller.ClientAbsent

	if value != nil {
		if value.client != nil {
			status |= controller.ClientCreated
		}
		if value.cancelFunc != nil {
			status |= controller.ClientStarted
		}
	}

	return status
}

func (s *ListenerController) Shutdown() error {
	s.state.listeners.Range(func(key string, value *Listener) bool {
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

func (s *ListenerController) ListenerRegister(name string, cfg *config.Listener) error {
	value, _ := s.state.listeners.Load(name)

	if status := s.statusRW(value); status != controller.ListenerAbsent {
		return statusErrorEq(name, status, controller.ListenerAbsent)
	}

	s.state.listeners.Store(name, &Listener{
		rw:         sync.RWMutex{},
		config:     cfg,
		cancelFunc: nil,
		listener:   nil,
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

		case *config.ListenerIpcHttp:
			listener, err := api.NewHttp(name, ipc.Port, s.plane)
			if err != nil {
				return err
			}
			value.listener = &ListenerMixinApi{
				listener: listener,
			}
			return nil
		}

	case *config.ListenerKindLsp:
		switch ipc := value.config.Ipc.(type) {

		case *config.ListenerIpcTcp:
			listener, err := lsp.NewTcp(name, ipc.Port, s.plane)
			if err != nil {
				return err
			}
			value.listener = &ListenerMixinLsp{
				listener: listener,
				clients:  xsync.NewMap[string, *ListenerLspClient](),
			}
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
		listener.Listener().Run(ctx)
		s.ListenerStop(listener.Listener().Name())
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

	switch value.config.Kind.(type) {
	case *config.ListenerKindLsp:
		listener := value.listener.(*ListenerMixinLsp)
		listener.clients.Range(func(key string, value *ListenerLspClient) bool {
			s.ListenerLspClientShutdown(listener.listener.Name(), key)
			return true
		})
	}

	cancel := value.cancelFunc
	value.cancelFunc = nil
	value.rw.Unlock()

	cancel()

	return nil
}

func (s *ListenerController) ListenerLspClientOwn(name string, client listener.ListenerLspClient) error {
	value, _ := s.state.listeners.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status&controller.ListenerStarted == 0 {
		return statusErrorFlag(name, status, controller.ListenerStarted)
	}

	switch kind := value.config.Kind.(type) {
	case *config.ListenerKindLsp:
	default:
		return fmt.Errorf("listener %v expected type %T, got %T", name, (*config.ListenerKindLsp)(nil), kind)
	}

	listener := value.listener.(*ListenerMixinLsp)
	clientv, _ := listener.clients.Load(client.Name())

	if status := s.lspClientStatus(clientv); status != controller.ClientAbsent {
		return clientStatusErrorEq(name, client.Name(), status, controller.ClientAbsent)
	}

	listener.clients.Store(client.Name(), &ListenerLspClient{
		client:     client,
		cancelFunc: nil,
	})

	return nil
}

func (s *ListenerController) ListenerLspClientStatus(name string, clientName string) controller.ClientStatus {
	value, _ := s.state.listeners.Load(name)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.ListenerStarted == 0 {
		return controller.ClientAbsent
	}

	switch value.config.Kind.(type) {
	case *config.ListenerKindLsp:
	default:
		return controller.ClientAbsent
	}

	listener := value.listener.(*ListenerMixinLsp)
	client, _ := listener.clients.Load(clientName)

	return s.lspClientStatus(client)
}

func (s *ListenerController) ListenerLspClientServe(name string, clientName string) error {
	value, _ := s.state.listeners.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status&controller.ListenerStarted == 0 {
		return statusErrorFlag(name, status, controller.ListenerStarted)
	}

	switch kind := value.config.Kind.(type) {
	case *config.ListenerKindLsp:
	default:
		return fmt.Errorf("listener %v expected type %T, got %T", name, (*config.ListenerKindLsp)(nil), kind)
	}

	listener := value.listener.(*ListenerMixinLsp)
	client, _ := listener.clients.Load(clientName)

	if status := s.lspClientStatus(client); status != controller.ClientCreated {
		return clientStatusErrorEq(name, clientName, status, controller.ClientCreated)
	}

	ctx, cancel := context.WithCancel(context.Background())

	client.cancelFunc = cancel

	go func() {
		client.client.Serve(ctx)
		s.ListenerLspClientShutdown(name, clientName)
	}()

	return nil
}

func (s *ListenerController) ListenerLspClientShutdown(name string, clientName string) error {
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

	switch kind := value.config.Kind.(type) {
	case *config.ListenerKindLsp:
	default:
		value.rw.Unlock()
		return fmt.Errorf("listener %v expected type %T, got %T", name, (*config.ListenerKindLsp)(nil), kind)
	}

	listener := value.listener.(*ListenerMixinLsp)
	client, _ := listener.clients.Load(clientName)

	if status := s.lspClientStatus(client); status&controller.ClientStarted == 0 {
		value.rw.Unlock()
		return clientStatusErrorFlag(name, clientName, controller.ClientStarted, status)
	}

	cancel := client.cancelFunc
	client.cancelFunc = nil

	listener.clients.Delete(clientName)

	value.rw.Unlock()

	cancel()

	return nil
}

func (s *ListenerController) ListenerGetConfig(name string) (*config.Listener, error) {
	value, _ := s.state.listeners.Load(name)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.ListenerRegistered == 0 {
		return nil, statusErrorFlag(name, status, controller.ListenerRegistered)
	}

	return value.config, nil
}

func (s *ListenerController) ListenerGet(name string) (*controller.ListenerData, error) {
	value, _ := s.state.listeners.Load(name)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	status := s.status(value)

	if status&controller.ListenerRegistered == 0 {
		return nil, statusErrorFlag(name, status, controller.ListenerRegistered)
	}

	config := value.config

	return &controller.ListenerData{
		Name:   name,
		Status: status,
		Config: config,
	}, nil
}

func (s *ListenerController) ListenerGetAll() []*controller.ListenerData {
	datas := make([]*controller.ListenerData, 0)

	s.state.listeners.Range(func(key string, value *Listener) bool {
		data, err := s.ListenerGet(key)

		if err == nil {
			datas = append(datas, data)
		}

		return true
	})

	return datas
}
