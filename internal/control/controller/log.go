package controller

import "mals/internal/log"

// NOTE: maybe creating log events is better
func (s *Controller) LogAdd(log log.Log) {
}

func (s *Controller) LogDelete(log log.Log) {
}

func (s *Controller) Debug(msg string, args ...any) {
	s.state.Logs.Range(func(key log.Log, value struct{}) bool {
		key.Debug(msg, args...)
		return true
	})
}

func (s *Controller) Info(msg string, args ...any) {
	s.state.Logs.Range(func(key log.Log, value struct{}) bool {
		key.Info(msg, args...)
		return true
	})
}

func (s *Controller) Warn(msg string, args ...any) {
	s.state.Logs.Range(func(key log.Log, value struct{}) bool {
		key.Warn(msg, args...)
		return true
	})
}

func (s *Controller) Error(msg string, args ...any) {
	s.state.Logs.Range(func(key log.Log, value struct{}) bool {
		key.Error(msg, args...)
		return true
	})
}
