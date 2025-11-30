package controller

import "mals/internal/control/event"

func (s *Controller) Shutdown() {
	s.bus.Publish(&event.EventShutdown{})
}
