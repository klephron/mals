package state

import "mals/internal/log"

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
