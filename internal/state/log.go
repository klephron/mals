package state

import (
	log "mals/internal/log/common"
)

func (s *State) LogAdd(log log.Log) {
	s.Logs = append(s.Logs, log)
}

func (s *State) LogContext() *log.LogGroup {
	return log.NewGroup(s.Logs)
}
