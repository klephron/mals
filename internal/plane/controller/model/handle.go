package model

import (
	"context"
	"fmt"
	"mals/internal/model"
	"mals/internal/model/openai"
	"mals/internal/plane/state"
	"mals/pkg/config"
)

func (s *ModelController) handleShutdown(t *TaskShutdown) {
	defer close(t.Result)

	s.state.Models.Range(func(key string, value *state.ModelValue) bool {
		ts := &TaskStop{TaskGeneric: NewTask(), Name: key}
		s.handleStop(ts)
		<-ts.Result

		td := &TaskDelete{TaskGeneric: NewTask(), Name: key}
		s.handleDelete(td)
		<-td.Result

		return true
	})

	t.Result <- nil
}

func (s *ModelController) handleRegister(t *TaskRegister) {
	defer close(t.Result)
	name := t.Config.Name

	if _, ok := s.state.Models.Load(name); ok {
		t.Result <- fmt.Errorf("model %v exists", name)
		return
	}

	switch settings := t.Config.Settings.(type) {
	case *config.ModelSettingsOpenAI:
		s.state.Models.Store(name, &state.ModelValue{
			Config:     t.Config,
			Model:      nil,
			CancelFunc: nil,
		})
		t.Result <- nil
	default:
		t.Result <- fmt.Errorf("unhandled model %T %v", settings, settings)
	}
}

func (s *ModelController) handleUnregister(t *TaskUnregister) {
	defer close(t.Result)
	name := t.Name

	value, ok := s.state.Models.Load(name)
	if !ok {
		t.Result <- fmt.Errorf("model %v does not exist", name)
		return
	}

	if value.Model != nil {
		t.Result <- fmt.Errorf("model %v is already created", name)
		return
	}

	s.state.Models.Delete(name)

	t.Result <- nil
}

func (s *ModelController) handleCreate(t *TaskCreate) {
	defer close(t.Result)
	name := t.Name

	value, ok := s.state.Models.Load(name)
	if !ok {
		t.Result <- fmt.Errorf("model %v does not exist", name)
		return
	}

	if value.Model != nil {
		t.Result <- fmt.Errorf("model %v is already created", name)
		return
	}

	switch settings := value.Config.Settings.(type) {
	case *config.ModelSettingsOpenAI:
		model, err := openai.New(name, openai.ModelOpenAISpec{
			MaxTokens:   settings.MaxTokens,
			Temperature: settings.Temperature,
		})
		if err != nil {
			t.Result <- err
			return
		}
		value.Model = model
		t.Result <- nil

	default:
		t.Result <- fmt.Errorf("unhandled model %T %v", settings, settings)
	}
}

func (s *ModelController) handleDelete(t *TaskDelete) {
	defer close(t.Result)
	name := t.Name

	value, ok := s.state.Models.Load(name)
	if !ok {
		t.Result <- fmt.Errorf("model %v does not exist", name)
		return
	}
	if value.Model == nil {
		t.Result <- fmt.Errorf("model %v is not created", name)
		return
	}
	if value.CancelFunc != nil {
		t.Result <- fmt.Errorf("model %v is running", name)
		return
	}

	value.Model = nil
	t.Result <- nil
}

func (s *ModelController) handleStart(t *TaskStart) {
	defer close(t.Result)
	name := t.Name

	value, ok := s.state.Models.Load(name)
	if !ok {
		t.Result <- fmt.Errorf("model %v does not exist", name)
		return
	}
	if value.Model == nil {
		t.Result <- fmt.Errorf("model %v is not created", name)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	value.CancelFunc = cancel

	go func(model model.Model) {
		model.Serve(ctx)
		// avoid race conditions
		s.Stop(model.Name())
	}(value.Model)

	t.Result <- nil
}

func (s *ModelController) handleStop(t *TaskStop) {
	defer close(t.Result)
	name := t.Name

	value, ok := s.state.Models.Load(name)
	if !ok {
		t.Result <- fmt.Errorf("model %v does not exist", name)
		return
	}
	if value.Model == nil {
		t.Result <- fmt.Errorf("model %v is not created", name)
		return
	}
	if value.CancelFunc == nil {
		t.Result <- fmt.Errorf("model %v is not running", name)
		return
	}

	value.CancelFunc()
	value.CancelFunc = nil

	t.Result <- nil
}
