package state

import (
	log "mals/internal/log"
)

func (s *StateImpl) LogAdd(log log.Log) {
	s.Logs.Store(log, struct{}{})
}

func (s *StateImpl) LogDelete(log log.Log) bool {
	_, loaded := s.Logs.LoadAndDelete(log)
	return loaded
}
