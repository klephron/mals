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

func (s *StateImpl) Debug(msg string, args ...any) {
	s.Logs.Range(func(key log.Log, value struct{}) bool {
		key.Debug(msg, args...)
		return true
	})
}

func (s *StateImpl) Info(msg string, args ...any) {
	s.Logs.Range(func(key log.Log, value struct{}) bool {
		key.Info(msg, args...)
		return true
	})
}

func (s *StateImpl) Warn(msg string, args ...any) {
	s.Logs.Range(func(key log.Log, value struct{}) bool {
		key.Warn(msg, args...)
		return true
	})
}

func (s *StateImpl) Error(msg string, args ...any) {
	s.Logs.Range(func(key log.Log, value struct{}) bool {
		key.Error(msg, args...)
		return true
	})
}
