package usage

import (
	"fmt"
	"mals/pkg/config"
)

func (s *UsageController) Shutdown() error {
	s.state.statusRW.RLock()
	cancel := s.state.statusCancel
	s.state.statusRW.RUnlock()

	cancel()

	return nil
}

func (s *UsageController) RegisterModel(config config.Model) error {
	name := config.Name

	if _, ok := s.state.configModels.Load(name); ok {
		return fmt.Errorf("config model %v exists", name)
	}

	s.state.configModels.Store(name, &config)

	return nil
}

func (s *UsageController) RegisterLsp(config config.Lsp) error {
	name := config.Name

	if _, ok := s.state.configLsps.Load(name); ok {
		return fmt.Errorf("config lsp %v exists", name)
	}

	s.state.configLsps.Store(name, &config)

	return nil
}

func (s *UsageController) RegisterUsage(config config.Usage) error {
	name := config.Name

	if _, ok := s.state.configUsages.Load(name); ok {
		return fmt.Errorf("config usage %v exists", name)
	}

	s.state.configUsages.Store(name, &config)

	return nil
}
