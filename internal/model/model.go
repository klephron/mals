package model

import (
	"context"

	"github.com/google/uuid"
)

type Task struct {
	Id                uuid.UUID `json:"id"`
	Text              string    `json:"text"`
	Schema            any       `json:"schema"`
	SchemaName        string    `json:"schema_name"`
	SchemaDescription string    `json:"schema_description"`
	SchemaStrict      bool      `json:"schema_strict"`
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
	Run(ctx context.Context) error
	Execute(ctx context.Context, task *Task) (string, error)
}
