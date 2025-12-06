package lifecycle

import "mals/internal/plane/event"

func (s *LifecycleController) handleLifecycleShutdown(_ *EventShutdown) {
	go func() {
		s.bus.Allcast(&event.EventShutdown{})
	}()
}

func (s *LifecycleController) handleLifecycleTerminate(_ *EventTerminate) {
	go func() {
		s.bus.Allcast(&event.EventTerminate{})
	}()
}
