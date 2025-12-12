package model

import (
	"context"

	"github.com/google/uuid"
)

type Result struct {
	Text  string
	Error error
}

type Task struct {
	id     uuid.UUID
	text   string
	result chan Result
}

func NewTask(text string) *Task {
	return &Task{id: uuid.New(), text: text, result: make(chan Result, 1)}
}

type Model interface {
	Name() string
	Kind() string
	Serve(ctx context.Context) error
	Submit(task Task) Result
}
