package scope

import (
	"fmt"
	"mals/pkg/config"
)

func (s *ScopeController) Shutdown() error {
	s.state.statusRW.RLock()
	cancel := s.state.statusCancel
	s.state.statusRW.RUnlock()

	cancel()

	return nil
}

func (s *ScopeController) RegisterModel(config config.Model) error {
	name := config.Name

	if _, ok := s.state.configModels.Load(name); ok {
		return fmt.Errorf("config model %v exists", name)
	}

	s.state.configModels.Store(name, &config)

	return nil
}

func (s *ScopeController) RegisterLsp(config config.Lsp) error {
	name := config.Name

	if _, ok := s.state.configLsps.Load(name); ok {
		return fmt.Errorf("config lsp %v exists", name)
	}

	s.state.configLsps.Store(name, &config)

	return nil
}
