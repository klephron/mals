package state

import (
	log "mals/internal/log/common"
)

func (s *State) LogAdd(log log.Log) {
	s.Logs.Store(log, struct{}{})
}

func (s *State) LogDelete(log log.Log) bool {
	_, loaded := s.Logs.LoadAndDelete(log)
	return loaded
}

func (s *State) Debug(msg string, args ...any) {
	s.Logs.Range(func(key log.Log, value struct{}) bool {
		key.Debug(msg, args...)
		return true
	})
}

func (s *State) Info(msg string, args ...any) {
	s.Logs.Range(func(key log.Log, value struct{}) bool {
		key.Info(msg, args...)
		return true
	})
}

func (s *State) Warn(msg string, args ...any) {
	s.Logs.Range(func(key log.Log, value struct{}) bool {
		key.Warn(msg, args...)
		return true
	})
}

func (s *State) Error(msg string, args ...any) {
	s.Logs.Range(func(key log.Log, value struct{}) bool {
		key.Error(msg, args...)
		return true
	})
}
