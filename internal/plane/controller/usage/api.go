package usage

import (
	"fmt"
	"mals/internal/usage"
	"mals/pkg/config"
	"slices"
)

func (s *UsageController) Shutdown() error {
	s.state.statusRW.RLock()
	cancel := s.state.statusCancel
	s.state.statusRW.RUnlock()

	cancel()

	return nil
}

func (s *UsageController) UsageRegister(cfg *config.Usage) error {
	name := cfg.Name

	if _, ok := s.state.usages.Load(name); ok {
		return fmt.Errorf("usage %v exists", name)
	}

	s.state.usages.Store(name, cfg)

	return nil
}

func (s *UsageController) UsageUnregister(name string) error {
	_, ok := s.state.usages.Load(name)
	if !ok {
		return fmt.Errorf("usage %v does not exist", name)
	}

	s.state.usages.Delete(name)

	return nil
}

func (s *UsageController) UsageGet(name string) (*config.Usage, error) {
	value, ok := s.state.usages.Load(name)
	if !ok {
		return nil, fmt.Errorf("usage %v does not exist", name)
	}

	return value, nil
}

func (s *UsageController) UsageGetFiltered(condition usage.ConditionFilter, event usage.EventFilter) []*config.Usage {
	usages := make([]*config.Usage, 0)

	s.state.usages.Range(func(key string, value *config.Usage) bool {
		usages = append(usages, value)
		return true
	})

	return usage.UsagesFilter(usages, condition, event)
}

func (s *UsageController) UsageGetFilteredClient(condition usage.ConditionFilter, event usage.EventFilter, listener string, client string) []*config.Usage {
	listenerConfig, err := s.plane.ListenerGetConfig(listener)
	if err != nil {
		s.plane.Warnf("%v", err)
		return nil
	}

	usages := make([]*config.Usage, 0)

	kindLsp, ok := listenerConfig.Kind.(*config.ListenerKindLsp)
	if !ok {
		return usages
	}

	s.state.usages.Range(func(key string, value *config.Usage) bool {
		if slices.Contains(kindLsp.Usages, value.Name) {
			usages = append(usages, value)
		}
		return true
	})

	return usage.UsagesFilter(usages, condition, event)
}

func (s *UsageController) UsageGetAll() []*config.Usage {
	usages := make([]*config.Usage, 0)

	s.state.usages.Range(func(key string, value *config.Usage) bool {
		usages = append(usages, value)
		return true
	})

	return usages
}
