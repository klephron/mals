package model

import (
	"context"
	"fmt"
	"mals/internal/client"
	"mals/internal/model"
	"mals/internal/model/metered"
	"mals/internal/model/openai"
	"mals/internal/plane/controller"
	"mals/pkg/config"
	"sync"

	"github.com/google/uuid"
)

func statusErrorEq(name string, actual controller.ModelStatus, expected controller.ModelStatus) error {
	return fmt.Errorf("model %v expected eq %v, got %v", name, expected, actual)
}

func statusErrorFlag(name string, actual controller.ModelStatus, expected controller.ModelStatus) error {
	return fmt.Errorf("model %v expected flag %v, got %v", name, expected, actual)
}

func (s *ModelController) status(value *ModelValue) controller.ModelStatus {
	status := controller.ModelAbsent

	if value != nil {
		status |= controller.ModelRegistered

		if value.model != nil {
			status |= controller.ModelCreated
		}
		if value.cancelFunc != nil {
			status |= controller.ModelStarted
		}
	}

	return status
}

func (s *ModelController) statusRW(value *ModelValue) controller.ModelStatus {
	status := controller.ModelAbsent

	if value != nil {
		status |= controller.ModelRegistered

		value.rw.RLock()

		if value.model != nil {
			status |= controller.ModelCreated
		}
		if value.cancelFunc != nil {
			status |= controller.ModelStarted
		}

		value.rw.RUnlock()
	}

	return status
}

func (s *ModelController) Shutdown() error {
	s.state.models.Range(func(key string, value *ModelValue) bool {
		s.ModelStop(key)
		s.ModelDelete(key)
		return true
	})

	s.state.statusRW.RLock()
	cancel := s.state.statusCancel
	s.state.statusRW.RUnlock()

	cancel()

	return nil
}

func (s *ModelController) ModelStatus(name string) controller.ModelStatus {
	value, _ := s.state.models.Load(name)
	return s.statusRW(value)
}

func (s *ModelController) ModelRegister(name string, cfg *config.Model) error {
	value, _ := s.state.models.Load(name)

	if status := s.statusRW(value); status != controller.ModelAbsent {
		return statusErrorEq(name, status, controller.ModelAbsent)
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

func (s *ModelController) ModelUnregister(name string) error {
	value, _ := s.state.models.Load(name)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status != controller.ModelRegistered {
		return statusErrorEq(name, status, controller.ModelRegistered)
	}

	s.state.models.Delete(name)

	return nil
}

func (s *ModelController) ModelCreate(name string) error {
	value, _ := s.state.models.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status != controller.ModelRegistered {
		return statusErrorEq(name, status, controller.ModelRegistered)
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

	if value.model != nil {
		value.model = metered.New(s.plane, value.model)
	}

	return nil
}

func (s *ModelController) ModelDelete(name string) error {
	value, _ := s.state.models.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status != controller.ModelRegistered|controller.ModelCreated {
		return statusErrorEq(name, status, controller.ModelRegistered|controller.ModelCreated)
	}

	value.model = nil
	value.queue = nil

	return nil
}

func (s *ModelController) ModelStart(name string) error {
	value, _ := s.state.models.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status != controller.ModelRegistered|controller.ModelCreated {
		return statusErrorEq(name, status, controller.ModelRegistered|controller.ModelCreated)
	}

	ctx, cancel := context.WithCancel(context.Background())

	value.cancelFunc = cancel

	model := value.model
	queue := value.queue

	go func() {
		go func() {
			err := model.Run(ctx)
			if err != nil {
				s.plane.Errorf("%v", err)
			}
			cancel()
		}()

		err := queue.serve(ctx, 1)
		if err != nil {
			s.plane.Errorf("%v", err)
		}

		s.ModelStop(model.Name())
	}()

	return nil
}

func (s *ModelController) ModelStop(name string) error {
	value, _ := s.state.models.Load(name)

	if value != nil {
		value.rw.Lock()
	}

	if status := s.status(value); status&controller.ModelStarted == 0 {
		if value != nil {
			value.rw.Unlock()
		}
		return statusErrorFlag(name, status, controller.ModelStarted)
	}

	cancel := value.cancelFunc
	value.cancelFunc = nil
	value.rw.Unlock()

	cancel()

	return nil
}

func (s *ModelController) TaskExecClient(modelName string, task *model.Task, client client.Client) (string, error) {
	value, _ := s.state.models.Load(modelName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.ModelCreated == 0 {
		return "", statusErrorFlag(modelName, status, controller.ModelCreated)
	}

	return value.queue.taskExecClient(task, client)
}

func (s *ModelController) TaskGetClient(modelName string, id uuid.UUID, client client.Client) (*model.Task, error) {
	value, _ := s.state.models.Load(modelName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.ModelCreated == 0 {
		return nil, statusErrorFlag(modelName, status, controller.ModelCreated)
	}

	return value.queue.taskGetClient(id, client)
}

func (s *ModelController) TaskGetAllClient(modelName string, client client.Client) ([]*model.Task, error) {
	value, _ := s.state.models.Load(modelName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.ModelCreated == 0 {
		return nil, statusErrorFlag(modelName, status, controller.ModelCreated)
	}

	return value.queue.taskGetAllClient(client), nil
}

func (s *ModelController) TaskCancelClient(modelName string, id uuid.UUID, client client.Client) (*model.Task, error) {
	value, _ := s.state.models.Load(modelName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.ModelCreated == 0 {
		return nil, statusErrorFlag(modelName, status, controller.ModelCreated)
	}

	return value.queue.taskCancelClient(id, client)
}

func (s *ModelController) TaskCancelAllClient(modelName string, client client.Client) ([]*model.Task, error) {
	value, _ := s.state.models.Load(modelName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.ModelCreated == 0 {
		return nil, statusErrorFlag(modelName, status, controller.ModelCreated)
	}

	return value.queue.taskCancelAllClient(client)
}
