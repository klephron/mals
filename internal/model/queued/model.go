package queued

import (
	"context"
	"mals/internal/model"
	"mals/internal/plane"
	"sync"

	"github.com/google/uuid"
)

type ModelQueued struct {
	model model.Model
	queue *taskQueue
	plane plane.Plane
}

func New(model model.Model, plane plane.Plane) *ModelQueued {
	s := &ModelQueued{
		model: model,
		queue: newTaskQueue(model),
		plane: plane,
	}
	return s
}

func (s *ModelQueued) Name() string {
	return s.model.Name()
}

func (s *ModelQueued) Kind() string {
	return s.model.Kind()
}

func (s *ModelQueued) Run(ctx context.Context, onReady func()) error {
	var wg sync.WaitGroup

	queueCtx, queueCancel := context.WithCancel(ctx)
	modelCtx, modelCancel := context.WithCancel(context.Background())

	wg.Go(func() {
		err := s.model.Run(modelCtx)
		if err != nil {
			s.plane.Errorf("%v", err)
		}
		queueCancel()
	})

	err := s.queue.serve(queueCtx, 1, func() { onReady() })

	modelCancel()
	wg.Wait()

	if err != nil {
		return s.plane.Errorf("%v", err)
	}

	return nil
}

func (s *ModelQueued) TaskExecClient(task *model.Task, clientName string) (string, error) {
	return s.queue.taskExecClient(task, clientName)
}

func (s *ModelQueued) TaskGet(id uuid.UUID) (*model.Task, error) {
	return s.queue.taskGet(id)
}

func (s *ModelQueued) TaskGetClient(id uuid.UUID, clientName string) (*model.Task, error) {
	return s.queue.taskGetClient(id, clientName)
}

func (s *ModelQueued) TaskGetAll() []*model.Task {
	return s.queue.taskGetAll()
}

func (s *ModelQueued) TaskGetAllClient(clientName string) []*model.Task {
	return s.queue.taskGetAllClient(clientName)
}

func (s *ModelQueued) TaskCancelClient(id uuid.UUID, clientName string) (*model.Task, error) {
	return s.queue.taskCancelClient(id, clientName)
}

func (s *ModelQueued) TaskCancelAllClient(clientName string) ([]*model.Task, error) {
	return s.queue.taskCancelAllClient(clientName)
}
