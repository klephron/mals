package model

import (
	"context"
	"fmt"
	"mals/internal/model"
	"mals/internal/model/openai"
	"mals/pkg/config"
)

func (s *ModelController) handleShutdown(t *TaskShutdown) {
	defer close(t.result)

	s.state.models.Range(func(key string, value *ModelValue) bool {
		ts := &TaskStop{TaskGeneric: NewTask(), name: key}
		s.handleStop(ts)
		<-ts.result

		td := &TaskDelete{TaskGeneric: NewTask(), name: key}
		s.handleDelete(td)
		<-td.result

		return true
	})

	t.result <- nil
}

func (s *ModelController) handleRegister(t *TaskRegister) {
	defer close(t.result)
	name := t.config.Name

	if _, ok := s.state.models.Load(name); ok {
		t.result <- fmt.Errorf("model %v exists", name)
		return
	}

	switch settings := t.config.Settings.(type) {
	case *config.ModelSettingsOpenAI:
		s.state.models.Store(name, &ModelValue{
			config:     t.config,
			model:      nil,
			cancelFunc: nil,
		})
		t.result <- nil
	default:
		t.result <- fmt.Errorf("unhandled model %T %v", settings, settings)
	}
}

func (s *ModelController) handleUnregister(t *TaskUnregister) {
	defer close(t.result)
	name := t.name

	value, ok := s.state.models.Load(name)
	if !ok {
		t.result <- fmt.Errorf("model %v does not exist", name)
		return
	}

	if value.model != nil {
		t.result <- fmt.Errorf("model %v is already created", name)
		return
	}

	s.state.models.Delete(name)

	t.result <- nil
}

func (s *ModelController) handleCreate(t *TaskCreate) {
	defer close(t.result)
	name := t.name

	value, ok := s.state.models.Load(name)
	if !ok {
		t.result <- fmt.Errorf("model %v does not exist", name)
		return
	}

	if value.model != nil {
		t.result <- fmt.Errorf("model %v is already created", name)
		return
	}

	switch settings := value.config.Settings.(type) {
	case *config.ModelSettingsOpenAI:
		model, err := openai.New(name, openai.ModelOpenAISpec{
			Url:         settings.Url,
			MaxTokens:   settings.MaxTokens,
			Temperature: settings.Temperature,
		})
		if err != nil {
			t.result <- err
			return
		}
		value.model = model
		t.result <- nil

	default:
		t.result <- fmt.Errorf("unhandled model %T %v", settings, settings)
	}
}

func (s *ModelController) handleDelete(t *TaskDelete) {
	defer close(t.result)
	name := t.name

	value, ok := s.state.models.Load(name)
	if !ok {
		t.result <- fmt.Errorf("model %v does not exist", name)
		return
	}
	if value.model == nil {
		t.result <- fmt.Errorf("model %v is not created", name)
		return
	}
	if value.cancelFunc != nil {
		t.result <- fmt.Errorf("model %v is running", name)
		return
	}

	value.model = nil
	t.result <- nil
}

func (s *ModelController) handleStart(t *TaskStart) {
	defer close(t.result)
	name := t.name

	value, ok := s.state.models.Load(name)
	if !ok {
		t.result <- fmt.Errorf("model %v does not exist", name)
		return
	}
	if value.model == nil {
		t.result <- fmt.Errorf("model %v is not created", name)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	value.cancelFunc = cancel

	go func(model model.Model) {
		model.Serve(ctx)
		// avoid race conditions
		s.Stop(model.Name())
	}(value.model)

	t.result <- nil
}

func (s *ModelController) handleStop(t *TaskStop) {
	defer close(t.result)
	name := t.name

	value, ok := s.state.models.Load(name)
	if !ok {
		t.result <- fmt.Errorf("model %v does not exist", name)
		return
	}
	if value.model == nil {
		t.result <- fmt.Errorf("model %v is not created", name)
		return
	}
	if value.cancelFunc == nil {
		t.result <- fmt.Errorf("model %v is not running", name)
		return
	}

	value.cancelFunc()
	value.cancelFunc = nil

	t.result <- nil
}
