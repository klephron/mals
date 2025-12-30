package metered

import (
	"context"
	"mals/internal/model"
	"time"
)

type ModelMetered struct {
	model.Model
}

func New(m model.Model) *ModelMetered {
	return &ModelMetered{Model: m}
}

func (s *ModelMetered) Execute(ctx context.Context, task *model.Task) (string, error) {
	start := time.Now()

	result, err := s.Model.Execute(task, ctx)
}
