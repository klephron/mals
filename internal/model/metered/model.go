package metered

import (
	"context"
	"mals/internal/model"
	"mals/internal/plane"
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
	parent model.Model
	plane  plane.Plane
}

func New(plane plane.Plane, model model.Model) *ModelMetered {
	s := &ModelMetered{
		parent: model,
		plane:  plane,
	}
	return s
}

func (s *ModelMetered) Name() string {
	return s.parent.Name()
}

func (s *ModelMetered) Kind() string {
	return s.parent.Kind()
}

func (s *ModelMetered) Run(ctx context.Context) error {
	return s.parent.Run(ctx)
}

func (s *ModelMetered) Execute(ctx context.Context, task *model.Task) (string, error) {
	start := time.Now()

	result, err := s.parent.Execute(ctx, task)

	duration := time.Since(start).Seconds()
	labels := prometheus.Labels{"model_name": s.Name(), "model_kind": s.Kind()}

	requestsProcessed.With(labels).Inc()
	executionDuration.With(labels).Observe(duration)

	return result, err
}
