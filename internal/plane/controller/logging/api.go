package logging

import (
	"mals/internal/log"
	"mals/internal/plane/state"
)

func (s *LogController) Debug(msg string, args ...any) {
	s.state.Logs.Range(func(key log.Log, value *state.LogValue) bool {
		if value.Enabled {
			key.Debug(msg, args...)
		}
		return true
	})
}

func (s *LogController) Info(msg string, args ...any) {
	s.state.Logs.Range(func(key log.Log, value *state.LogValue) bool {
		if value.Enabled {
			key.Info(msg, args...)
		}
		return true
	})
}

func (s *LogController) Warn(msg string, args ...any) {
	s.state.Logs.Range(func(key log.Log, value *state.LogValue) bool {
		if value.Enabled {
			key.Warn(msg, args...)
		}
		return true
	})
}

func (s *LogController) Error(msg string, args ...any) {
	s.state.Logs.Range(func(key log.Log, value *state.LogValue) bool {
		if value.Enabled {
			key.Error(msg, args...)
		}
		return true
	})
}
