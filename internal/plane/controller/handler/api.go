package handler

import (
	"fmt"
	"mals/pkg/config"
)

func (s *HandlerController) Register(cfg *config.Handler) error {
	name := cfg.Name

	if _, ok := s.state.handlers.Load(name); ok {
		return fmt.Errorf("usage %v exists", name)
	}

	s.state.handlers.Store(name, cfg)

	return nil
}

func (s *HandlerController) Unregister(name string) error {
	_, ok := s.state.handlers.Load(name)
	if !ok {
		return fmt.Errorf("usage %v does not exist", name)
	}

	s.state.handlers.Delete(name)

	return nil
}

func (s *HandlerController) Get(name string) (*config.Handler, error) {
	value, ok := s.state.handlers.Load(name)
	if !ok {
		return nil, fmt.Errorf("usage %v does not exist", name)
	}

	return value, nil
}

func (s *HandlerController) GetAll() []*config.Handler {
	usages := make([]*config.Handler, 0)

	s.state.handlers.Range(func(key string, value *config.Handler) bool {
		usages = append(usages, value)
		return true
	})

	return usages
}
