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
			fmt.Println("TODO: gracefully shutdown")
			return
		default:
			s.Warn(fmt.Sprintf("received event %T %v", e, e))
			fmt.Printf("received event %T %v\n", e, e)
		}
	}
}
