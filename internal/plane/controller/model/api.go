package model

import (
	"context"
	"fmt"
	"mals/internal/model/openai"
	"mals/pkg/config"
	"sync"
)

func (s *ModelController) Shutdown() error {
	s.state.models.Range(func(key string, value *ModelValue) bool {
		s.Stop(key)
		s.Delete(key)
		return true
	})

	s.state.statusRW.RLock()
	cancel := s.state.statusCancel
	s.state.statusRW.RUnlock()

	cancel()

	return nil
}

func (s *ModelController) Register(cfg config.Model) error {
	name := cfg.Name

	if _, ok := s.state.models.Load(name); ok {
		return fmt.Errorf("model %v exists", name)
	}

	switch settings := cfg.Settings.(type) {
	case *config.ModelSettingsOpenAI:
		s.state.models.Store(name, &ModelValue{
			rw:         sync.RWMutex{},
			config:     cfg,
			model:      nil,
			cancelFunc: nil,
		})
	default:
		return fmt.Errorf("unhandled model %T %v", settings, settings)
	}

	return nil
}

func (s *ModelController) Unregister(name string) error {
	value, ok := s.state.models.Load(name)
	if !ok {
		return fmt.Errorf("model %v does not exist", name)
	}

	value.rw.RLock()

	if value.model != nil {
		value.rw.RUnlock()
		return fmt.Errorf("model %v is already created", name)
	}

	s.state.models.Delete(name)
	value.rw.RUnlock()

	return nil
}

func (s *ModelController) Create(name string) error {
	value, ok := s.state.models.Load(name)
	if !ok {
		return fmt.Errorf("model %v does not exist", name)
	}

	value.rw.Lock()
	defer value.rw.Unlock()

	if value.model != nil {
		return fmt.Errorf("model %v is already created", name)
	}

	switch settings := value.config.Settings.(type) {
	case *config.ModelSettingsOpenAI:
		model, err := openai.New(name, openai.ModelOpenAISpec{
			Url:         settings.Url,
			MaxTokens:   settings.MaxTokens,
			Temperature: settings.Temperature,
		})
		if err != nil {
			return err
		}
		value.model = model

	default:
		return fmt.Errorf("unhandled model %T %v", settings, settings)
	}

	return nil
}

func (s *ModelController) Delete(name string) error {
	value, ok := s.state.models.Load(name)
	if !ok {
		return fmt.Errorf("model %v does not exist", name)
	}

	value.rw.Lock()

	if value.model == nil {
		value.rw.Unlock()
		return fmt.Errorf("model %v is not created", name)
	}
	if value.cancelFunc != nil {
		value.rw.Unlock()
		return fmt.Errorf("model %v is running", name)
	}

	value.model = nil
	value.rw.Unlock()

	return nil
}

func (s *ModelController) Start(name string) error {
	value, ok := s.state.models.Load(name)
	if !ok {
		return fmt.Errorf("model %v does not exist", name)
	}

	value.rw.Lock()

	if value.model == nil {
		value.rw.Unlock()
		return fmt.Errorf("model %v is not created", name)
	}

	ctx, cancel := context.WithCancel(context.Background())

	value.cancelFunc = cancel
	model := value.model
	value.rw.Unlock()

	go func() {
		model.Serve(ctx)
		s.Stop(model.Name())
	}()

	return nil
}

func (s *ModelController) Stop(name string) error {
	value, ok := s.state.models.Load(name)
	if !ok {
		return fmt.Errorf("model %v does not exist", name)
	}

	value.rw.Lock()

	if value.model == nil {
		value.rw.Unlock()
		return fmt.Errorf("model %v is not created", name)
	}
	if value.cancelFunc == nil {
		value.rw.Unlock()
		return fmt.Errorf("model %v is not running", name)
	}

	cancel := value.cancelFunc
	value.cancelFunc = nil
	value.rw.Unlock()

	cancel()

	return nil
}

// func (s *ModelController) TaskExecClient(model string, task model.Task, client client.Client) model.Result {
// }

// func (s *ModelController) TaskGetClient(model string, id uuid.UUID, client client.Client) (Task, error) {

// }

// func (s *ModelController) TaskGetAllClient(model string, client client.Client) ([]Task, error) {

// }

// func (s *ModelController) TaskGetAllClientName(model string, client string) ([]Task, error) {

// }

// func (s *ModelController) TaskDeleteClient(model string, id uuid.UUID, client client.Client) (Task, error) {

// }

// func (s *ModelController) TaskDeleteAllClient(model string, id uuid.UUID, client client.Client) ([]Task, error) {

// }
