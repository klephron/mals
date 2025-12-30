package metered

import (
	"context"
	"mals/internal/model"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "model_requests_processed_total",
			Help: "Total number of requests processed by model",
		},
		[]string{"model_name", "model_kind"},
	)

	executionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "model_execution_duration_seconds",
			Help:    "Execution duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"model_name", "model_kind"},
	)
)

type ModelMetered struct {
	model.Model
}

func New(m model.Model) *ModelMetered {
	return &ModelMetered{Model: m}
}

func (s *ModelMetered) Execute(ctx context.Context, task *model.Task) (string, error) {
	start := time.Now()

	result, err := s.Model.Execute(ctx, task)

	duration := time.Since(start).Seconds()
	labels := prometheus.Labels{"model_name": s.Name(), "model_kind": s.Kind()}

	requestsProcessed.With(labels).Inc()
	executionDuration.With(labels).Observe(duration)

	return result, err
}
