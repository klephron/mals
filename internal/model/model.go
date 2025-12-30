package model

import (
	"context"

	"github.com/google/uuid"
)

type Task struct {
	Id                uuid.UUID
	Text              string
	Schema            any
	SchemaName        string
	SchemaDescription string
	SchemaStrict      bool
}

func NewTask(text string, schema any, schemaName string, schemaDescription string) *Task {
	return &Task{
		Id:                uuid.New(),
		Text:              text,
		Schema:            schema,
		SchemaName:        schemaName,
		SchemaDescription: schemaDescription,
		SchemaStrict:      true,
	}
}

type Model interface {
	Name() string
	Kind() string
	Serve(ctx context.Context) error
	Execute(ctx context.Context, task *Task) (string, error)
}
