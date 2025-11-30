package state

type StateLog struct {
	Enabled bool
}

func NewStateLog(enabled bool) *StateLog {
	return &StateLog{
		Enabled: enabled,
	}
}
