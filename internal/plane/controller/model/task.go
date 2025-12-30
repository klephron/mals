package model

import (
	"context"
	"fmt"
	"mals/internal/client"
	"mals/internal/model"
	"sync"

	"github.com/google/uuid"
	"github.com/puzpuzpuz/xsync/v4"
)

type TaskResult struct {
	text  string
	error error
}

type TaskRequest struct {
	client client.Client
	task   *model.Task
	result chan TaskResult
	ctx    context.Context
	cancel context.CancelFunc
}

type TaskQueue struct {
	rw     sync.RWMutex
	mapped *xsync.Map[uuid.UUID, *TaskRequest]
	queued chan *TaskRequest
	ctx    context.Context
	model  model.Model
}

func newTaskQueue(model model.Model) *TaskQueue {
	return &TaskQueue{
		rw:     sync.RWMutex{},
		mapped: xsync.NewMap[uuid.UUID, *TaskRequest](),
		queued: nil,
		ctx:    nil,
		model:  model,
	}
}

func (s *TaskQueue) serve(ctx context.Context, workers int, onReady func()) error {
	var wg sync.WaitGroup

	if workers <= 0 {
		return fmt.Errorf("%T %v worker cnt must be > 0", s, s.model.Name())
	}

	s.rw.Lock()
	if s.ctx != nil {
		s.rw.Unlock()
		return fmt.Errorf("%T %v queue is serving", s, s.model.Name())
	}

	s.queued = make(chan *TaskRequest)
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

func (s *TaskQueue) worker() {
	s.rw.RLock()
	queued := s.queued
	s.rw.RUnlock()

	for request := range queued {
		select {
		case <-request.ctx.Done():
			request.result <- TaskResult{error: fmt.Errorf("task %v cancelled", request.task.Id)}
		default:
			text, error := s.model.Execute(request.ctx, request.task)
			request.result <- TaskResult{text: text, error: error}
		}
	}
}

func (s *TaskQueue) taskExecClient(task *model.Task, client client.Client) (string, error) {
	s.rw.RLock()

	if s.ctx == nil {
		s.rw.RUnlock()
		return "", fmt.Errorf("%T %v queue is not serving", s, s.model.Name())
	}

	taskCtx, taskCancel := context.WithCancel(s.ctx)

	request := &TaskRequest{
		client: client,
		task:   task,
		result: make(chan TaskResult, 1),
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

func (s *TaskQueue) taskCancelClient(id uuid.UUID, client client.Client) (*model.Task, error) {
	request, ok := s.mapped.Load(id)
	if !ok || request.client != client {
		return nil, fmt.Errorf("task %v not found", id)
	}
	request.cancel()
	return request.task, nil
}

func (s *TaskQueue) taskCancelAllClient(client client.Client) ([]*model.Task, error) {
	ids := make([]*model.Task, 0)
	s.mapped.Range(func(key uuid.UUID, value *TaskRequest) bool {
		if value.client != client {
			return true
		}
		value.cancel()
		ids = append(ids, value.task)
		return true
	})
	return ids, nil
}

func (s *TaskQueue) taskGetClient(id uuid.UUID, client client.Client) (*model.Task, error) {
	request, ok := s.mapped.Load(id)
	if !ok || request.client != client {
		return nil, fmt.Errorf("task %v not found", id)
	}
	return request.task, nil
}

func (s *TaskQueue) taskGet(id uuid.UUID) (*model.Task, error) {
	request, ok := s.mapped.Load(id)
	if !ok {
		return nil, fmt.Errorf("task %v not found", id)
	}
	return request.task, nil
}

func (s *TaskQueue) taskGetAllClient(client client.Client) []*model.Task {
	tasks := make([]*model.Task, 0)
	s.mapped.Range(func(key uuid.UUID, value *TaskRequest) bool {
		if value.client != client {
			return true
		}
		tasks = append(tasks, value.task)
		return true
	})
	return tasks
}

func (s *TaskQueue) taskGetAll() []*model.Task {
	tasks := make([]*model.Task, 0)
	s.mapped.Range(func(key uuid.UUID, value *TaskRequest) bool {
		tasks = append(tasks, value.task)
		return true
	})
	return tasks
}
