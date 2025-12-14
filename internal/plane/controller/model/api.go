package model

import (
	"context"
	"fmt"
	"mals/internal/client"
	"mals/internal/model"
	"mals/internal/model/openai"
	"mals/pkg/config"
	"sync"

	"github.com/google/uuid"
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
			queue:      nil,
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
		spec := openai.ModelOpenAISpec{
			Url:         settings.Url,
			MaxTokens:   settings.MaxTokens,
			Temperature: settings.Temperature,
		}
		model, err := openai.New(name, spec, s.plane)
		if err != nil {
			return err
		}
		value.model = model
		value.queue = newTaskQueue(value.model)

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
	value.queue = nil
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
	queue := value.queue

	value.rw.Unlock()

	go func() {
		go func() {
			err := model.Serve(ctx)
			if err != nil {
				s.plane.Log().Errorf("%v", err)
			}
			cancel()
		}()

		err := queue.serve(ctx, 1)
		if err != nil {
			s.plane.Log().Errorf("%v", err)
		}

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

func (s *ModelController) TaskExecClient(modelName string, task *model.Task, client client.Client) model.Result {
	value, ok := s.state.models.Load(modelName)
	if !ok {
		return model.Result{Error: fmt.Errorf("model %v does not exist", modelName)}
	}

	value.rw.RLock()
	defer value.rw.RUnlock()

	if value.model == nil {
		return model.Result{Error: fmt.Errorf("model %v is not created", modelName)}
	}
	if value.cancelFunc == nil {
		return model.Result{Error: fmt.Errorf("model %v is not running", modelName)}
	}

	return value.queue.taskExecClient(task, client)
}

func (s *ModelController) TaskGetClient(modelName string, id uuid.UUID, client client.Client) (*model.Task, error) {
	value, ok := s.state.models.Load(modelName)
	if !ok {
		return nil, fmt.Errorf("model %v does not exist", modelName)
	}

	value.rw.RLock()
	defer value.rw.RUnlock()

	if value.model == nil {
		return nil, fmt.Errorf("model %v is not created", modelName)
	}

	return value.queue.taskGetClient(id, client)
}

func (s *ModelController) TaskGetAllClient(modelName string, client client.Client) ([]*model.Task, error) {
	value, ok := s.state.models.Load(modelName)
	if !ok {
		return nil, fmt.Errorf("model %v does not exist", modelName)
	}

	value.rw.RLock()
	defer value.rw.RUnlock()

	if value.model == nil {
		return nil, fmt.Errorf("model %v is not created", modelName)
	}

	return value.queue.taskGetAllClient(client), nil
}

func (s *ModelController) TaskCancelClient(modelName string, id uuid.UUID, client client.Client) (*model.Task, error) {
	value, ok := s.state.models.Load(modelName)
	if !ok {
		return nil, fmt.Errorf("model %v does not exist", modelName)
	}

	value.rw.RLock()
	defer value.rw.RUnlock()

	if value.model == nil {
		return nil, fmt.Errorf("model %v is not created", modelName)
	}

	return value.queue.taskCancelClient(id, client)
}

func (s *ModelController) TaskCancelAllClient(modelName string, client client.Client) ([]*model.Task, error) {
	value, ok := s.state.models.Load(modelName)
	if !ok {
		return nil, fmt.Errorf("model %v does not exist", modelName)
	}

	value.rw.RLock()
	defer value.rw.RUnlock()

	if value.model == nil {
		return nil, fmt.Errorf("model %v is not created", modelName)
	}

	return value.queue.taskCancelAllClient(client)
}
