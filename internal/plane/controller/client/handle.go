package client

import (
	"mals/internal/client"
	"mals/internal/plane/event"
	"mals/internal/plane/state"
)

func (s *ClientController) handleShutdown(_ *event.EventShutdown) {
	s.state.Clients.Range(func(key client.Client, value *state.ClientValue) bool {
		s.handleStop(&TaskStop{TaskGeneric: NewTaskSingle(), Client: key})
		s.handleDelete(&TaskDelete{TaskGeneric: NewTaskSingle(), Client: key, Listener: value.Listener})
		return true
	})
}

func (s *ClientController) handleTerminate(_ *event.EventTerminate) {
	s.state.Clients.Range(func(key client.Client, value *state.ClientValue) bool {
		s.state.Clients.Delete(key)
		return true
	})
}

func (s *ClientController) handleOwn(t *TaskOwn) {

}

func (s *ClientController) handleDelete(t *TaskDelete) {

}

func (s *ClientController) handleStart(t *TaskStart) {

}

func (s *ClientController) handleStop(t *TaskStop) {

}
