package model

import (
	"context"
	"mals/pkg/core"

	"github.com/google/uuid"
)

type Task struct {
	Id           uuid.UUID           `json:"id"`
	Messages     []core.ModelMessage `json:"messages,omitempty"`
	Schema       core.ModelSchema    `json:"schema"`
	SchemaStrict bool                `json:"schema_strict"`
}

func NewTaskWPrompt(prompt string, schema core.ModelSchema) *Task {
	return &Task{
		Id:           uuid.New(),
		Messages:     []core.ModelMessage{{Role: core.ModelRoleUser, Content: prompt}},
		Schema:       schema,
		SchemaStrict: true,
	}
}

func NewTaskWMessages(messages []core.ModelMessage, schema core.ModelSchema) *Task {
	return &Task{
		Id:           uuid.New(),
		Messages:     messages,
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
