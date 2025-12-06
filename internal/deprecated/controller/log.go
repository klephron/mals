package controller

import (
	"mals/internal/control/event"
	"mals/internal/log"
)

func (s *Controller) LogAdd(log log.Log) {
	s.bus.Publish(&event.EventLogAdd{Log: log})
}

func (s *Controller) LogDelete(log log.Log) {
	s.bus.Publish(&event.EventLogDelete{Log: log})
}

func (s *Controller) LogStart(log log.Log) {
	s.bus.Publish(&event.EventLogStart{Log: log})
}

func (s *Controller) LogStop(log log.Log) {
	s.bus.Publish(&event.EventLogStop{Log: log})
}
