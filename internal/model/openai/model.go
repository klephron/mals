package openai

import (
	"context"
	"fmt"
	"mals/internal/model"
)

type ModelOpenAISpec struct {
	Url         string
	MaxTokens   int
	Temperature float32
}

type ModelOpenAI struct {
	model.Model
	name string
	spec ModelOpenAISpec
}

func New(name string, spec ModelOpenAISpec) (*ModelOpenAI, error) {
	return &ModelOpenAI{
		name: name,
		spec: spec,
	}, nil
}

func (s *ModelOpenAI) Name() string {
	return s.name
}

func (s *ModelOpenAI) Kind() string {
	return Kind()
}

func (s *ModelOpenAI) Serve(ctx context.Context) error {
	return nil
}

func (s *ModelOpenAI) Submit(task model.Task) model.Result {
	return model.Result{
		Text:  "",
		Error: fmt.Errorf("connection is undefined"),
	}
}
