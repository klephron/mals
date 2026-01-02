package lsp

import (
	"fmt"
	"mals/internal/plane/controller"
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
