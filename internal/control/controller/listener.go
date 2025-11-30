package controller

import (
	"mals/internal/control/event"
	"mals/internal/listener"
)

func (s *Controller) ListenerAdd(Listener listener.Listener) {
	s.bus.Publish(&event.EventListenerAdd{Listener: Listener})
}

func (s *Controller) ListenerDelete(Listener listener.Listener) {
	s.bus.Publish(&event.EventListenerDelete{Listener: Listener})
}

func (s *Controller) ListenerStart(Listener listener.Listener) {
	s.bus.Publish(&event.EventListenerStart{Listener: Listener})
}

func (s *Controller) ListenerStop(Listener listener.Listener) {
	s.bus.Publish(&event.EventListenerStop{Listener: Listener})
}
