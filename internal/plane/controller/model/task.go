package model

import (
	"context"
	"fmt"
	"mals/internal/model"
	"sync"

	"github.com/google/uuid"
	"github.com/puzpuzpuz/xsync/v4"
)

type TaskRequest struct {
	task   model.Task
	result chan model.Result
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

func (s *TaskQueue) serve(ctx context.Context, workers int) error {
	var wg sync.WaitGroup

	if workers <= 0 {
		return fmt.Errorf("%T %v worker cnt must be > 0", s, s.model.Name())
	}

	s.rw.Lock()
	if s.ctx != nil {
		s.rw.Unlock()
		return fmt.Errorf("%T %v is serving", s, s.model.Name())
	}

	s.queued = make(chan *TaskRequest)
	s.ctx = ctx
	s.rw.Unlock()

	for range workers {
		wg.Go(func() {
			s.worker()
		})
	}

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
			request.result <- model.Result{Error: fmt.Errorf("task %v cancelled", request.task.Id)}
		default:
			request.result <- s.model.Execute(request.task, request.ctx)
		}
	}
}

func (s *TaskQueue) taskExec(task model.Task) model.Result {
	s.rw.RLock()

	if s.ctx == nil {
		s.rw.RUnlock()
		return model.Result{Error: fmt.Errorf("%T %v is not serving", s, s.model.Name())}
	}

	taskCtx, taskCancel := context.WithCancel(s.ctx)

	request := &TaskRequest{
		task:   task,
		result: make(chan model.Result, 1),
		ctx:    taskCtx,
		cancel: taskCancel,
	}

	s.mapped.Store(task.Id, request)
	s.queued <- request
	s.rw.RUnlock()

	result := <-request.result
	s.mapped.Delete(task.Id)
	return result
}

func (s *TaskQueue) taskCancel(id uuid.UUID) error {
	request, ok := s.mapped.Load(id)
	if !ok {
		return fmt.Errorf("task %v not found", id)
	}
	request.cancel()
	return nil
}

func (s *TaskQueue) taskGetAll() []model.Task {
	ids := make([]model.Task, 0)
	s.mapped.Range(func(key uuid.UUID, value *TaskRequest) bool {
		ids = append(ids, value.task)
		return true
	})
	return ids
}
