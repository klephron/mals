package scheduler

import (
	"context"
	"fmt"
	"mals/internal/control/event"
	"mals/internal/control/state"
	"mals/internal/listener"
	"mals/internal/log"
)

func (s *Scheduler) EventLoop() {
	for e := range s.ch {
		switch e := e.(type) {
		case *event.EventLogAdd:
			s.HandleLogAdd(e)
		case *event.EventLogDelete:
			s.HandleLogDelete(e)
		case *event.EventLogStart:
			s.HandleLogStart(e)
		case *event.EventLogStop:
			s.HandleLogStop(e)

		case *event.EventListenerAdd:
			s.HandleListenerAdd(e)
		case *event.EventListenerDelete:
			s.HandleListenerDelete(e)
		case *event.EventListenerStart:
			s.HandleListenerStart(e)
		case *event.EventListenerStop:
			s.HandleListenerStop(e)

		case *event.EventShutdown:
			s.HandleShutdown(e)
		case *event.EventTerminate:
			return

		default:
			s.Warn(fmt.Sprintf("%T %v unhandled", e, e))
		}
	}
}

func (s *Scheduler) HandleShutdown(e *event.EventShutdown) {
	s.state.Listeners.Range(func(key listener.Listener, value *state.StateListener) bool {
		s.HandleListenerDelete(&event.EventListenerDelete{Listener: key})
		return true
	})
	s.state.Logs.Range(func(key log.Log, value *state.StateLog) bool {
		s.HandleLogDelete(&event.EventLogDelete{Log: key})
		return true
	})
	go func() {
		s.bus.Publish(&event.EventTerminate{})
	}()
}

func (s *Scheduler) HandleListenerAdd(e *event.EventListenerAdd) {
	_, exist := s.state.Listeners.Load(e.Listener)
	if exist {
		return
	}
	state := state.NewStateListener()
	s.state.Listeners.Store(e.Listener, state)
	s.Info(fmt.Sprintf("HandleListenerAdd: %T %v", e.Listener, e.Listener))
}

func (s *Scheduler) HandleListenerDelete(e *event.EventListenerDelete) {
	_, exist := s.state.Listeners.Load(e.Listener)
	if !exist {
		return
	}
	s.HandleListenerStop(&event.EventListenerStop{Listener: e.Listener})
	s.state.Listeners.Delete(e.Listener)
	s.Info(fmt.Sprintf("HandleListenerDelete: %T %v", e.Listener, e.Listener))
}

func (s *Scheduler) HandleListenerStart(e *event.EventListenerStart) {
	state, exist := s.state.Listeners.Load(e.Listener)
	if !exist || state.CancelFunc != nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	state.CancelFunc = cancel
	go func() {
		e.Listener.Listen(ctx)
	}()
	s.Info(fmt.Sprintf("HandleListenerStart: %T %v", e.Listener, e.Listener))
}

func (s *Scheduler) HandleListenerStop(e *event.EventListenerStop) {
	state, exist := s.state.Listeners.Load(e.Listener)
	if !exist {
		return
	}
	if state.CancelFunc != nil {
		state.CancelFunc()
		state.CancelFunc = nil
	}
	s.Info(fmt.Sprintf("HandleListenerStop: %T %v", e.Listener, e.Listener))
}

func (s *Scheduler) HandleLogAdd(e *event.EventLogAdd) {
	_, exist := s.state.Logs.Load(e.Log)
	if exist {
		return
	}
	state := state.NewStateLog(false)
	s.state.Logs.Store(e.Log, state)
	s.Info(fmt.Sprintf("HandleLogAdd: %T %v", e.Log, e.Log))
}

func (s *Scheduler) HandleLogDelete(e *event.EventLogDelete) {
	_, exist := s.state.Logs.LoadAndDelete(e.Log)
	if !exist {
		return
	}
	e.Log.Close()
	s.Info(fmt.Sprintf("HandleLogDelete: %T %v", e.Log, e.Log))
}

func (s *Scheduler) HandleLogStart(e *event.EventLogStart) {
	state, exist := s.state.Logs.Load(e.Log)
	if !exist {
		return
	}
	state.Enabled = true
	s.Info(fmt.Sprintf("HandleLogStart: %T %v", e.Log, e.Log))
}

func (s *Scheduler) HandleLogStop(e *event.EventLogStop) {
	state, exist := s.state.Logs.Load(e.Log)
	if !exist {
		return
	}
	state.Enabled = false
	s.Info(fmt.Sprintf("HandleLogStop: %T %v", e.Log, e.Log))
}
