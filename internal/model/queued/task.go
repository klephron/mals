package queued

import (
	"context"
	"fmt"
	"mals/internal/model"
	"sync"

	"github.com/google/uuid"
	"github.com/puzpuzpuz/xsync/v4"
)

type taskResult struct {
	text  string
	error error
}

type taskRequest struct {
	task   *model.Task
	client string
	result chan taskResult
	ctx    context.Context
	cancel context.CancelFunc
}

type taskQueue struct {
	rw     sync.RWMutex
	mapped *xsync.Map[uuid.UUID, *taskRequest]
	queued chan *taskRequest
	ctx    context.Context
	model  model.Model
}

func newTaskQueue(model model.Model) *taskQueue {
	return &taskQueue{
		rw:     sync.RWMutex{},
		mapped: xsync.NewMap[uuid.UUID, *taskRequest](),
		queued: nil,
		ctx:    nil,
		model:  model,
	}
}

func (s *taskQueue) serve(ctx context.Context, workers int, onReady func()) error {
	var wg sync.WaitGroup

	if workers <= 0 {
		return fmt.Errorf("%T %v worker cnt must be > 0", s, s.model.Name())
	}

	s.rw.Lock()
	if s.ctx != nil {
		s.rw.Unlock()
		return fmt.Errorf("%T %v queue is serving", s, s.model.Name())
	}

	s.queued = make(chan *taskRequest)
	s.ctx = ctx
	s.rw.Unlock()

	for range workers {
		wg.Go(func() {
			s.worker()
		})
	}

	onReady()

	<-ctx.Done()

	s.rw.Lock()
	close(s.queued)
	s.queued = nil
	s.ctx = nil
	s.rw.Unlock()

	wg.Wait()

	return nil
}

func (s *taskQueue) worker() {
	s.rw.RLock()
	queued := s.queued
	s.rw.RUnlock()

	for request := range queued {
		select {
		case <-request.ctx.Done():
			request.result <- taskResult{error: fmt.Errorf("task %v cancelled", request.task.Id)}
		default:
			text, error := s.model.Execute(request.ctx, request.task)
			request.result <- taskResult{text: text, error: error}
		}
	}
}

func (s *taskQueue) taskExecClient(task *model.Task, clientName string) (string, error) {
	s.rw.RLock()

	if s.ctx == nil {
		s.rw.RUnlock()
		return "", fmt.Errorf("%T %v queue is not serving", s, s.model.Name())
	}

	taskCtx, taskCancel := context.WithCancel(s.ctx)

	request := &taskRequest{
		client: clientName,
		task:   task,
		result: make(chan taskResult, 1),
		ctx:    taskCtx,
		cancel: taskCancel,
	}

	s.mapped.Store(task.Id, request)
	s.queued <- request
	s.rw.RUnlock()

	result := <-request.result
	s.mapped.Delete(task.Id)
	return result.text, result.error
}

func (s *taskQueue) taskCancelClient(id uuid.UUID, clientName string) (*model.Task, error) {
	request, ok := s.mapped.Load(id)
	if !ok || request.client != clientName {
		return nil, fmt.Errorf("task %v not found", id)
	}
	request.cancel()
	return request.task, nil
}

func (s *taskQueue) taskCancelAllClient(clientName string) ([]*model.Task, error) {
	ids := make([]*model.Task, 0)
	s.mapped.Range(func(key uuid.UUID, value *taskRequest) bool {
		if value.client != clientName {
			return true
		}
		value.cancel()
		ids = append(ids, value.task)
		return true
	})
	return ids, nil
}

func (s *taskQueue) taskGetClient(id uuid.UUID, clientName string) (*model.Task, error) {
	request, ok := s.mapped.Load(id)
	if !ok || request.client != clientName {
		return nil, fmt.Errorf("task %v not found", id)
	}
	return request.task, nil
}

func (s *taskQueue) taskGet(id uuid.UUID) (*model.Task, error) {
	request, ok := s.mapped.Load(id)
	if !ok {
		return nil, fmt.Errorf("task %v not found", id)
	}
	return request.task, nil
}

func (s *taskQueue) taskGetAllClient(clientName string) []*model.Task {
	tasks := make([]*model.Task, 0)
	s.mapped.Range(func(key uuid.UUID, value *taskRequest) bool {
		if value.client != clientName {
			return true
		}
		tasks = append(tasks, value.task)
		return true
	})
	return tasks
}

func (s *taskQueue) taskGetAll() []*model.Task {
	tasks := make([]*model.Task, 0)
	s.mapped.Range(func(key uuid.UUID, value *taskRequest) bool {
		tasks = append(tasks, value.task)
		return true
	})
	return tasks
}
