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
	Id   uuid.UUID
	Text string
}

func NewTask(text string) *Task {
	return &Task{Id: uuid.New(), Text: text}
}

type Model interface {
	Name() string
	Kind() string
	Serve(ctx context.Context) error
	Execute(task Task, ctx context.Context) Result
}
