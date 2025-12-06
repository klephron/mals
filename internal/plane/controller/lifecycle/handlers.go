package lifecycle

import "mals/internal/plane/event"

func (s *LifecycleController) handleShutdown(e *event.EventShutdown) {

}

func (s *LifecycleController) handleTerminate(_ *event.EventTerminate) {
}
