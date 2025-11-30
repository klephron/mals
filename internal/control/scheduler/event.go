package scheduler

import (
	"fmt"
	"mals/internal/control/event"
)

func (s *Scheduler) EventLoop() {
	for e := range s.ch {
		switch e := e.(type) {
		case *event.EventShutdown:
			s.Warn("TODO: gracefully shutdown")
			return
		default:
			s.Warn(fmt.Sprintf("received event %v of type %T", e, e))
		}
	}
}
