package lsp

import (
	"context"
	"fmt"
	"mals/internal/lsp/protocol"
	"mals/internal/lsp/server/stdio"
	"mals/internal/plane/controller"
	"mals/pkg/config"
	"sync"
)

func statusErrorEq(name string, actual controller.LspStatus, expected controller.LspStatus) error {
	return fmt.Errorf("model %v expected eq %v, got %v", name, expected, actual)
}

func statusErrorFlag(name string, actual controller.LspStatus, expected controller.LspStatus) error {
	return fmt.Errorf("model %v expected flag %v, got %v", name, expected, actual)
}

func (s *LspController) status(value *LspValue) controller.LspStatus {
	status := controller.LspAbsent

	if value != nil {
		status |= controller.LspRegistered

		if value.lsp != nil {
			status |= controller.LspCreated
		}
		if value.cancelFunc != nil {
			status |= controller.LspStarted
		}
	}

	return status
}

func (s *LspController) statusRW(value *LspValue) controller.LspStatus {
	status := controller.LspAbsent

	if value != nil {
		status |= controller.LspRegistered

		value.rw.RLock()

		if value.lsp != nil {
			status |= controller.LspCreated
		}
		if value.cancelFunc != nil {
			status |= controller.LspStarted
		}

		value.rw.RUnlock()
	}

	return status
}

func (s *LspController) Shutdown() error {
	s.state.lsps.Range(func(key string, value *LspValue) bool {
		s.LspStop(key)
		s.LspDelete(key)
		return true
	})

	s.state.statusRW.RLock()
	cancel := s.state.statusCancel
	s.state.statusRW.RUnlock()

	cancel()

	return nil
}

func (s *LspController) LspStatus(name string) controller.LspStatus {
	value, _ := s.state.lsps.Load(name)
	return s.statusRW(value)
}

func (s *LspController) LspRegister(name string, cfg *config.Lsp) error {
	value, _ := s.state.lsps.Load(name)

	if status := s.statusRW(value); status != controller.LspAbsent {
		return statusErrorEq(name, status, controller.LspAbsent)
	}

	switch settings := cfg.Settings.(type) {
	case *config.LspSettingsStdio:
		s.state.lsps.Store(name, &LspValue{
			rw:         sync.RWMutex{},
			config:     cfg,
			lsp:        nil,
			cancelFunc: nil,
		})
	default:
		return fmt.Errorf("unhandled lsp %T %v", settings, settings)
	}

	return nil
}

func (s *LspController) LspUnregister(name string) error {
	value, _ := s.state.lsps.Load(name)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status != controller.LspRegistered {
		return statusErrorEq(name, status, controller.LspRegistered)
	}

	s.state.lsps.Delete(name)

	return nil
}

func (s *LspController) LspCreate(name string) error {
	value, _ := s.state.lsps.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status != controller.LspRegistered {
		return statusErrorEq(name, status, controller.LspRegistered)
	}

	switch settings := value.config.Settings.(type) {
	case *config.LspSettingsStdio:
		value.lsp = stdio.New(name, settings, s.plane)

	default:
		return fmt.Errorf("unhandled model %T %v", settings, settings)
	}

	return nil
}

func (s *LspController) LspDelete(name string) error {
	value, _ := s.state.lsps.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status != controller.LspRegistered|controller.LspCreated {
		return statusErrorEq(name, status, controller.LspRegistered|controller.LspCreated)
	}

	value.lsp = nil

	return nil
}

func (s *LspController) LspStart(name string) error {
	value, _ := s.state.lsps.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status != controller.LspRegistered|controller.LspCreated {
		return statusErrorEq(name, status, controller.LspRegistered|controller.LspCreated)
	}

	ctx, cancel := context.WithCancel(context.Background())

	value.cancelFunc = cancel

	lsp := value.lsp

	var wgReady sync.WaitGroup

	wgReady.Add(1)

	go func() {
		err := lsp.Run(ctx, func() { wgReady.Done() })
		if err != nil {
			s.plane.Errorf("%v", err)
		}
		cancel()

		s.LspStop(lsp.Name())
	}()

	wgReady.Wait()

	return nil
}

func (s *LspController) LspStop(name string) error {
	value, _ := s.state.lsps.Load(name)

	if value != nil {
		value.rw.Lock()
	}

	if status := s.status(value); status&controller.LspStarted == 0 {
		if value != nil {
			value.rw.Unlock()
		}
		return statusErrorFlag(name, status, controller.LspStarted)
	}

	cancel := value.cancelFunc
	value.cancelFunc = nil
	value.rw.Unlock()

	cancel()

	return nil
}

func (s *LspController) LspGetCapabilities(lspName string) (*protocol.ServerCapabilities, error) {
	value, _ := s.state.lsps.Load(lspName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.LspStarted == 0 {
		return nil, statusErrorFlag(lspName, status, controller.LspStarted)
	}

	return value.lsp.Capabilities()
}

func (s *LspController) LspGetInfo(lspName string) (*protocol.ServerInfo, error) {
	value, _ := s.state.lsps.Load(lspName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.LspStarted == 0 {
		return nil, statusErrorFlag(lspName, status, controller.LspStarted)
	}

	return value.lsp.Info()
}

func (s *LspController) LspGet(lspName string) (*controller.LspData, error) {
	value, _ := s.state.lsps.Load(lspName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	status := s.status(value)

	if status&controller.LspRegistered == 0 {
		return nil, statusErrorFlag(lspName, status, controller.LspRegistered)
	}

	config := value.config

	capabilities, _ := s.LspGetCapabilities(lspName)
	info, _ := s.LspGetInfo(lspName)

	return &controller.LspData{
		Name:         lspName,
		Status:       status,
		Config:       config,
		Capabilities: capabilities,
		Info:         info,
	}, nil
}

func (s *LspController) LspGetAll() []*controller.LspData {
	datas := make([]*controller.LspData, 0)

	s.state.lsps.Range(func(key string, value *LspValue) bool {
		data, err := s.LspGet(key)

		if err == nil {
			datas = append(datas, data)
		}

		return true
	})

	return datas
}

func (s *LspController) EventInitialize(lspName string, params *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	value, _ := s.state.lsps.Load(lspName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.LspStarted == 0 {
		return nil, statusErrorFlag(lspName, status, controller.LspStarted)
	}

	return value.lsp.Initialize(params)
}

func (s *LspController) EventInitialized(lspName string, params *protocol.InitializedParams) error {
	value, _ := s.state.lsps.Load(lspName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.LspStarted == 0 {
		return statusErrorFlag(lspName, status, controller.LspStarted)
	}

	return value.lsp.Initialized(params)
}

func (s *LspController) EventTextDocumentDidOpen(lspName string, params *protocol.DidOpenTextDocumentParams) error {
	value, _ := s.state.lsps.Load(lspName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.LspStarted == 0 {
		return statusErrorFlag(lspName, status, controller.LspStarted)
	}

	return value.lsp.TextDocumentDidOpen(params)
}

func (s *LspController) EventTextDocumentDidChange(lspName string, params *protocol.DidChangeTextDocumentParams) error {
	value, _ := s.state.lsps.Load(lspName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.LspStarted == 0 {
		return statusErrorFlag(lspName, status, controller.LspStarted)
	}

	return value.lsp.TextDocumentDidChange(params)
}

func (s *LspController) EventTextDocumentDidClose(lspName string, params *protocol.DidCloseTextDocumentParams) error {
	value, _ := s.state.lsps.Load(lspName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.LspStarted == 0 {
		return statusErrorFlag(lspName, status, controller.LspStarted)
	}

	return value.lsp.TextDocumentDidClose(params)
}

func (s *LspController) EventTextDocumentCompletion(lspName string, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	value, _ := s.state.lsps.Load(lspName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.LspStarted == 0 {
		return nil, statusErrorFlag(lspName, status, controller.LspStarted)
	}

	return value.lsp.TextDocumentCompletion(params)
}

func (s *LspController) EventShutdown(lspName string) error {
	value, _ := s.state.lsps.Load(lspName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.LspStarted == 0 {
		return statusErrorFlag(lspName, status, controller.LspStarted)
	}

	return value.lsp.Shutdown()
}
