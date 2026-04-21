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

		if value.mixin != nil {
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

		if value.mixin != nil {
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

func (s *ListenerController) Status(name string) controller.ListenerStatus {
	value, _ := s.state.listeners.Load(name)
	return s.statusRW(value)
}

func (s *ListenerController) Register(name string, cfg *config.Listener) error {
	value, _ := s.state.listeners.Load(name)

	if status := s.statusRW(value); status != controller.ListenerAbsent {
		return statusErrorEq(name, status, controller.ListenerAbsent)
	}

	s.state.listeners.Store(name, &Listener{
		rw:         sync.RWMutex{},
		config:     cfg,
		cancelFunc: nil,
		mixin:      nil,
	})

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

	switch value.config.Protocol.(type) {
	case *config.ListenerProtocolApi:
		switch ipc := value.config.Ipc.(type) {

		case *config.ListenerIpcTcp:
			if ipc.Port == nil {
				return fmt.Errorf("listener %v port is not set", value.config)
			}
			listener, err := api.NewHttp(name, int(*ipc.Port), s.plane)
			if err != nil {
				return err
			}
			value.mixin = &ListenerMixinApi{
				listener: listener,
			}
			return nil
		}

	case *config.ListenerProtocolLsp:
		switch ipc := value.config.Ipc.(type) {

		case *config.ListenerIpcTcp:
			listener, err := lsp.NewTcp(name, int(*ipc.Port), s.plane)
			if err != nil {
				return err
			}
			value.mixin = &ListenerMixinLsp{
				listener: listener,
				clients:  xsync.NewMap[string, *ListenerLspClient](),
			}
			return nil
		}
	}

	return fmt.Errorf("unhandled listener %v protocol=%v ipc=%v", value.config, value.config.Protocol.ListenerProtocolKind(), value.config.Ipc.ListenerIpcKind())
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

	listener := value.mixin.Listener()

	value.mixin = nil

	s.plane.Infof("%T: %T %v deleted", s, listener, listener.Name())

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
	listener := value.mixin.Listener()

	go func() {
		listener.Run(ctx)
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

	switch value.config.Protocol.(type) {
	case *config.ListenerProtocolLsp:
		mixinLsp := value.mixin.(*ListenerMixinLsp)

		mixinLsp.clients.Range(func(key string, value *ListenerLspClient) bool {
			s.lspClientShutdown(mixinLsp, key)
			return true
		})
	}

	listener := value.mixin.Listener()

	s.plane.Debugf("%T: %T %v to be stopped", s, listener, listener.Name())

	cancel := value.cancelFunc
	value.cancelFunc = nil
	value.rw.Unlock()

	cancel()

	s.plane.Infof("%T: %T %v stopped", s, listener, listener.Name())

	return nil
}

func (s *ListenerController) LspClientOwn(name string, client listener.ListenerLspClient) error {
	value, _ := s.state.listeners.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status&controller.ListenerStarted == 0 {
		return statusErrorFlag(name, status, controller.ListenerStarted)
	}

	switch kind := value.config.Protocol.(type) {
	case *config.ListenerProtocolLsp:
	default:
		return fmt.Errorf("listener %v expected type %T, got %T", name, (*config.ListenerProtocolLsp)(nil), kind)
	}

	mixinLsp := value.mixin.(*ListenerMixinLsp)
	clientValue, _ := mixinLsp.clients.Load(client.Name())

	if status := s.lspClientStatus(clientValue); status != controller.ClientAbsent {
		return clientStatusErrorEq(name, client.Name(), status, controller.ClientAbsent)
	}

	mixinLsp.clients.Store(client.Name(), &ListenerLspClient{
		client:     client,
		cancelFunc: nil,
	})

	return nil
}

func (s *ListenerController) LspClientStatus(name string, clientName string) controller.ClientStatus {
	value, _ := s.state.listeners.Load(name)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.ListenerStarted == 0 {
		return controller.ClientAbsent
	}

	switch value.config.Protocol.(type) {
	case *config.ListenerProtocolLsp:
	default:
		return controller.ClientAbsent
	}

	mixinLsp := value.mixin.(*ListenerMixinLsp)
	clientValue, _ := mixinLsp.clients.Load(clientName)

	return s.lspClientStatus(clientValue)
}

func (s *ListenerController) LspClientServe(name string, clientName string) error {
	value, _ := s.state.listeners.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status&controller.ListenerStarted == 0 {
		return statusErrorFlag(name, status, controller.ListenerStarted)
	}

	switch kind := value.config.Protocol.(type) {
	case *config.ListenerProtocolLsp:
	default:
		return fmt.Errorf("listener %v expected type %T, got %T", name, (*config.ListenerProtocolLsp)(nil), kind)
	}

	mixinLsp := value.mixin.(*ListenerMixinLsp)
	clientValue, _ := mixinLsp.clients.Load(clientName)

	if status := s.lspClientStatus(clientValue); status != controller.ClientCreated {
		return clientStatusErrorEq(name, clientName, status, controller.ClientCreated)
	}

	ctx, cancel := context.WithCancel(context.Background())

	clientValue.cancelFunc = cancel

	go func() {
		clientValue.client.Serve(ctx)
		s.lspClientShutdown(mixinLsp, clientName)
	}()

	return nil
}

func (s *ListenerController) lspClientShutdown(mixinLsp *ListenerMixinLsp, clientName string) error {
	name := mixinLsp.listener.Name()
	value, _ := mixinLsp.clients.Load(clientName)

	if status := s.lspClientStatus(value); status&controller.ClientStarted == 0 {
		return clientStatusErrorFlag(name, clientName, controller.ClientStarted, status)
	}

	cancel := value.cancelFunc
	value.cancelFunc = nil

	mixinLsp.clients.Delete(clientName)

	s.plane.Infof("%T: %T %v deleted", s, value.client, clientName)

	s.plane.Debugf("%T: %T %v to be stopped", s, value.client, clientName)

	cancel()

	s.plane.Infof("%T: %T %v stopped", s, value.client, clientName)

	return nil
}

func (s *ListenerController) LspClientShutdown(name string, clientName string) error {
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

	switch kind := value.config.Protocol.(type) {
	case *config.ListenerProtocolLsp:
	default:
		value.rw.Unlock()
		return fmt.Errorf("listener %v expected type %T, got %T", name, (*config.ListenerProtocolLsp)(nil), kind)
	}

	mixinLsp := value.mixin.(*ListenerMixinLsp)

	defer value.rw.Unlock()

	return s.lspClientShutdown(mixinLsp, clientName)
}

func (s *ListenerController) GetConfig(name string) (*config.Listener, error) {
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

func (s *ListenerController) Get(name string) (*controller.ListenerData, error) {
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

func (s *ListenerController) GetAll() []*controller.ListenerData {
	datas := make([]*controller.ListenerData, 0)

	s.state.listeners.Range(func(key string, value *Listener) bool {
		data, err := s.Get(key)

		if err == nil {
			datas = append(datas, data)
		}

		return true
	})

	return datas
}
