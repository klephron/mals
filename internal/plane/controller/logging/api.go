package logging

import (
	"fmt"
	"mals/internal/plane/state"
)

func (s *LogController) Debugf(format string, a ...any) {
	s.state.Logs.Range(func(key string, value *state.LogValue) bool {
		if value.Log != nil && value.Enabled {
			value.Log.Debug(fmt.Sprintf(format, a...))
		}
		return true
	})
}

func (s *LogController) Infof(format string, a ...any) {
	s.state.Logs.Range(func(key string, value *state.LogValue) bool {
		if value.Log != nil && value.Enabled {
			value.Log.Info(format, a...)
		}
		return true
	})
}

func (s *LogController) Warnf(format string, a ...any) {
	s.state.Logs.Range(func(key string, value *state.LogValue) bool {
		if value.Log != nil && value.Enabled {
			value.Log.Warn(format, a...)
		}
		return true
	})
}

func (s *LogController) Errorf(format string, a ...any) {
	s.state.Logs.Range(func(key string, value *state.LogValue) bool {
		if value.Log != nil && value.Enabled {
			value.Log.Error(format, a...)
		}
		return true
	})
}
