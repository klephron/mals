package model

import (
	"context"
	"fmt"
	"mals/internal/model"
	"mals/internal/model/metered"
	"mals/internal/model/openai"
	"mals/internal/model/queued"
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

func (s *ModelController) status(value *stateModel) controller.ModelStatus {
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

func (s *ModelController) statusRW(value *stateModel) controller.ModelStatus {
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

func (s *ModelController) Status(name string) controller.ModelStatus {
	value, _ := s.state.models.Load(name)
	return s.statusRW(value)
}

func (s *ModelController) Register(name string, cfg *config.Model) error {
	value, _ := s.state.models.Load(name)

	if status := s.statusRW(value); status != controller.ModelAbsent {
		return statusErrorEq(name, status, controller.ModelAbsent)
	}

	switch settings := cfg.Api.(type) {
	case *config.ModelApiOpenai:
		s.state.models.Store(name, &stateModel{
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

func (s *ModelController) Create(name string) error {
	value, _ := s.state.models.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status != controller.ModelRegistered {
		return statusErrorEq(name, status, controller.ModelRegistered)
	}

	var model model.Model

	switch api := value.config.Api.(type) {
	case *config.ModelApiOpenai:
		model = openai.New(name, api, s.plane)

	default:
		return fmt.Errorf("unhandled model %T %v", api, api)
	}

	model = metered.New(model, s.plane)
	value.model = queued.New(model, s.plane)

	return nil
}

func (s *ModelController) Delete(name string) error {
	value, _ := s.state.models.Load(name)

	if value != nil {
		value.rw.Lock()
		defer value.rw.Unlock()
	}

	if status := s.status(value); status != controller.ModelRegistered|controller.ModelCreated {
		return statusErrorEq(name, status, controller.ModelRegistered|controller.ModelCreated)
	}

	value.model = nil

	return nil
}

func (s *ModelController) Start(name string) error {
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

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		model.Run(ctx, func() { wg.Done() })
		s.Stop(model.Name())
	}()

	wg.Wait()

	return nil
}

func (s *ModelController) Stop(name string) error {
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

func (s *ModelController) Get(name string) (*controller.ModelData, error) {
	value, _ := s.state.models.Load(name)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	status := s.status(value)

	if status&controller.ModelRegistered == 0 {
		return nil, statusErrorFlag(name, status, controller.ModelRegistered)
	}

	config := value.config

	tasks, _ := s.TaskGetAll(name)

	return &controller.ModelData{
		Name:   name,
		Status: status,
		Config: config,
		Tasks:  tasks,
	}, nil
}

func (s *ModelController) GetAll() []*controller.ModelData {
	datas := make([]*controller.ModelData, 0)

	s.state.models.Range(func(key string, value *stateModel) bool {
		data, err := s.Get(key)

		if err == nil {
			datas = append(datas, data)
		}

		return true
	})

	return datas
}

func (s *ModelController) TaskExecClient(modelName string, task *model.Task, clientName string) (string, error) {
	value, _ := s.state.models.Load(modelName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.ModelCreated == 0 {
		return "", statusErrorFlag(modelName, status, controller.ModelCreated)
	}

	return value.model.TaskExecClient(task, clientName)
}

func (s *ModelController) TaskGet(modelName string, id uuid.UUID) (*model.Task, error) {
	value, _ := s.state.models.Load(modelName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.ModelCreated == 0 {
		return nil, statusErrorFlag(modelName, status, controller.ModelCreated)
	}

	return value.model.TaskGet(id)
}

func (s *ModelController) TaskGetClient(modelName string, id uuid.UUID, clientName string) (*model.Task, error) {
	value, _ := s.state.models.Load(modelName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.ModelCreated == 0 {
		return nil, statusErrorFlag(modelName, status, controller.ModelCreated)
	}

	return value.model.TaskGetClient(id, clientName)
}

func (s *ModelController) TaskGetAll(modelName string) ([]*model.Task, error) {
	value, _ := s.state.models.Load(modelName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.ModelCreated == 0 {
		return nil, statusErrorFlag(modelName, status, controller.ModelCreated)
	}

	return value.model.TaskGetAll(), nil
}

func (s *ModelController) TaskGetAllClient(modelName string, clientName string) ([]*model.Task, error) {
	value, _ := s.state.models.Load(modelName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.ModelCreated == 0 {
		return nil, statusErrorFlag(modelName, status, controller.ModelCreated)
	}

	return value.model.TaskGetAllClient(clientName), nil
}

func (s *ModelController) TaskCancelClient(modelName string, id uuid.UUID, clientName string) (*model.Task, error) {
	value, _ := s.state.models.Load(modelName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.ModelCreated == 0 {
		return nil, statusErrorFlag(modelName, status, controller.ModelCreated)
	}

	return value.model.TaskCancelClient(id, clientName)
}

func (s *ModelController) TaskCancelAllClient(modelName string, clientName string) ([]*model.Task, error) {
	value, _ := s.state.models.Load(modelName)

	if value != nil {
		value.rw.RLock()
		defer value.rw.RUnlock()
	}

	if status := s.status(value); status&controller.ModelCreated == 0 {
		return nil, statusErrorFlag(modelName, status, controller.ModelCreated)
	}

	return value.model.TaskCancelAllClient(clientName)
}
