package state

import "context"

type StateListener struct {
	CancelFunc context.CancelFunc
}

func NewStateListener() *StateListener {
	return &StateListener{
		CancelFunc: nil,
	}
}
