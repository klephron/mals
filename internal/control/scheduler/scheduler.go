package scheduler

import (
	"mals/internal/control/event"
	"mals/internal/control/state"
	"mals/internal/control/util"
)

type Scheduler struct {
	bus   *event.EventBus
	state *state.State
	ch    <-chan event.Event
}

func New(state *state.State, bus *event.EventBus) *Scheduler {
	return &Scheduler{
		state: state,
		bus:   bus,
	}
}

func (s *Scheduler) Subscribe() {
	if s.ch != nil {
		return
	}
	s.ch = s.bus.Subscribe()
}

func (s *Scheduler) ServeSubscribed() {
	if s.ch == nil {
		panic("must can subscribe before serve")
	}

	defer func() {
		s.bus.Unsubscribe(s.ch)
		s.ch = nil
	}()

	s.EventLoop()
}

func (s *Scheduler) Debug(msg string, args ...any) {
	util.Debug(s.state, msg, args...)
}

func (s *Scheduler) Info(msg string, args ...any) {
	util.Info(s.state, msg, args...)
}

func (s *Scheduler) Warn(msg string, args ...any) {
	util.Warn(s.state, msg, args...)
}

func (s *Scheduler) Error(msg string, args ...any) {
	util.Error(s.state, msg, args...)
}
