package scheduler

import (
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
	go func(s *Scheduler) {
		s.state.Listeners.Range(func(key listener.Listener, value struct{}) bool {
			s.bus.Publish(&event.EventListenerDelete{Listener: key})
			return true
		})
		s.state.Logs.Range(func(key log.Log, value *state.StateLog) bool {
			s.bus.Publish(&event.EventLogDelete{Log: key})
			return true
		})
		s.bus.Publish(&event.EventTerminate{})
	}(s)
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
