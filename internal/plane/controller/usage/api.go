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

func (s *UsageController) Register(cfg config.Usage) error {
	name := cfg.Name

	if _, ok := s.state.usages.Load(name); ok {
		return fmt.Errorf("usage %v exists", name)
	}

	s.state.usages.Store(name, &cfg)

	return nil
}

func (s *UsageController) Unregister(name string) error {
	_, ok := s.state.usages.Load(name)
	if !ok {
		return fmt.Errorf("usage %v does not exist", name)
	}

	s.state.usages.Delete(name)

	return nil
}

func (s *UsageController) GetAll() []*config.Usage {
	usages := make([]*config.Usage, 0)

	s.state.usages.Range(func(key string, value *config.Usage) bool {
		usages = append(usages, value)
		return true
	})

	return usages
}
