package state

import "context"

type ClientState struct {
	Cancel context.CancelFunc
}

func NewClientState(cancel context.CancelFunc) *ClientState {
	return &ClientState{
		Cancel: cancel,
	}
}
