package lifecycle

import "mals/internal/plane/event"

func (s *LifecycleController) Shutdown() {
	s.bus.Allcast(&event.EventShutdown{})
}

func (s *LifecycleController) Terminate() {
	s.bus.Allcast(&event.EventTerminate{})
}
