package model

import (
	"context"
	"mals/pkg/config"

	"github.com/google/uuid"
)

type Task struct {
	Id           uuid.UUID              `json:"id"`
	Text         string                 `json:"text"`
	Schema       config.StepModelSchema `json:"schema"`
	SchemaStrict bool                   `json:"schema_strict"`
}

func NewTask(text string, schema config.StepModelSchema) *Task {
	return &Task{
		Id:           uuid.New(),
		Text:         text,
		Schema:       schema,
		SchemaStrict: true,
	}
}

type Model interface {
	Name() string
	Kind() string
	Run(ctx context.Context) error
	Execute(ctx context.Context, task *Task) (string, error)
}
