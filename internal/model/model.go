package model

import (
	"context"
	"mals/pkg/core"

	"github.com/google/uuid"
)

type Task struct {
	Id         uuid.UUID            `json:"id"`
	Messages   []core.ModelMessage  `json:"messages,omitempty"`
	Parameters core.ModelParameters `json:"parameters"`
}

func NewTaskWPrompt(prompt string, parameters core.ModelParameters) *Task {
	return &Task{
		Id:         uuid.New(),
		Messages:   []core.ModelMessage{{Role: core.ModelRoleUser, Content: prompt}},
		Parameters: parameters,
	}
}

func NewTaskWMessages(messages []core.ModelMessage, parameters core.ModelParameters) *Task {
	return &Task{
		Id:         uuid.New(),
		Messages:   messages,
		Parameters: parameters,
	}
}

type Model interface {
	Name() string
	Kind() string
	Run(ctx context.Context) error
	Execute(ctx context.Context, task *Task) (string, error)
}
