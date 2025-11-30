package scheduler

import (
	"fmt"
	"mals/internal/control/event"
	"mals/internal/control/state"
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

func (s *Scheduler) EventLoop() {
	for event := range s.ch {
		switch event := event.(type) {
		default:
			fmt.Printf("received event %v of type %T", event, event)
		}
	}
}
