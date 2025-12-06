package logging

import "mals/internal/plane/event"

func (s *LogController) handleShutdown(e *event.EventShutdown) {

}

func (s *LogController) handleTerminate(_ *event.EventTerminate) {
}
