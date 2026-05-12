package lsp

import (
	"context"
	"fmt"
	"mals/internal/lsp/server/stdio"
	"mals/internal/plane/controller"
	"mals/pkg/config"
	"mals/third_party/lsp"
	"sync"
)

func statusErrorEq(name string, actual controller.LspStatus, expected controller.LspStatus) error {
	return fmt.Errorf("model %v expected eq %v, got %v", name, expected, actual)
}

func statusErrorFlag(name string, actual controller.LspStatus, expected controller.LspStatus) error {
	return fmt.Errorf("model %v expected flag %v, got %v", name, expected, actual)
}

func (s *LspController) status(value *stateLsp) controller.LspStatus {
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

func (s *LspController) statusRW(value *stateLsp) controller.LspStatus {
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

func (s *LspController) Status(name string) controller.LspStatus {
	value, _ := s.state.lsps.Load(name)
	return s.statusRW(value)
}

func (s *LspController) Register(name string, cfg *config.Lsp) error {
	value, _ := s.state.lsps.Load(name)

	if status := s.statusRW(value); status != controller.LspAbsent {
		return statusErrorEq(name, status, controller.LspAbsent)
	}

	switch settings := cfg.Api.(type) {
	case *config.LspApiStdio:
		s.state.lsps.Store(name, &stateLsp{
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

func (s *LspController) Unregister(name string) error {
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

func (s *LspController) Create(name string) error {
	value, _ := s.state.lsps.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status != controller.LspRegistered {
		return statusErrorEq(name, status, controller.LspRegistered)
	}

	switch settings := value.config.Api.(type) {
	case *config.LspApiStdio:
		value.lsp = stdio.New(name, settings, s.plane)

	default:
		return fmt.Errorf("unhandled model %T %v", settings, settings)
	}

	return nil
}

func (s *LspController) Delete(name string) error {
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

func (s *LspController) Start(name string) error {
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

		s.Stop(lsp.Name())
	}()

	wgReady.Wait()

	return nil
}

func (s *LspController) Stop(name string) error {
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

func (s *LspController) GetCapabilities(lspName string) (*lsp.ServerCapabilities, error) {
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

func (s *LspController) GetInfo(lspName string) (*lsp.ServerInfo, error) {
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

func (s *LspController) Get(lspName string) (*controller.LspData, error) {
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

	capabilities, _ := s.GetCapabilities(lspName)
	info, _ := s.GetInfo(lspName)

	return &controller.LspData{
		Name:         lspName,
		Status:       status,
		Config:       config,
		Capabilities: capabilities,
		Info:         info,
	}, nil
}

func (s *LspController) GetAll() []*controller.LspData {
	datas := make([]*controller.LspData, 0)

	s.state.lsps.Range(func(key string, value *stateLsp) bool {
		data, err := s.Get(key)

		if err == nil {
			datas = append(datas, data)
		}

		return true
	})

	return datas
}

func (s *LspController) Initialize(lspName string, params *lsp.InitializeParams) (*lsp.InitializeResult, error) {
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

func (s *LspController) Initialized(lspName string, params *lsp.InitializedParams) error {
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

func (s *LspController) TextDocumentDidOpen(lspName string, params *lsp.DidOpenTextDocumentParams) error {
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

func (s *LspController) TextDocumentDidChange(lspName string, params *lsp.DidChangeTextDocumentParams) error {
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

func (s *LspController) TextDocumentDidClose(lspName string, params *lsp.DidCloseTextDocumentParams) error {
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

func (s *LspController) TextDocumentCompletion(lspName string, params *lsp.CompletionParams) (*lsp.CompletionList, error) {
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

func (s *LspController) Shutdown(lspName string) error {
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
