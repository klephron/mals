package util

import (
	"mals/internal/control/state"
	"mals/internal/log"
)

func Debug(s *state.State, msg string, args ...any) {
	s.Logs.Range(func(key log.Log, value *state.StateLog) bool {
		if value.Enabled {
			key.Debug(msg, args...)
		}
		return true
	})
}

func Info(s *state.State, msg string, args ...any) {
	s.Logs.Range(func(key log.Log, value *state.StateLog) bool {
		if value.Enabled {
			key.Info(msg, args...)
		}
		return true
	})
}

func Warn(s *state.State, msg string, args ...any) {
	s.Logs.Range(func(key log.Log, value *state.StateLog) bool {
		if value.Enabled {
			key.Warn(msg, args...)
		}
		return true
	})
}

func Error(s *state.State, msg string, args ...any) {
	s.Logs.Range(func(key log.Log, value *state.StateLog) bool {
		if value.Enabled {
			key.Error(msg, args...)
		}
		return true
	})
}
