package client

import (
	"mals/internal/client"
	"mals/internal/listener"
)

type Task interface {
	task()
}

type TaskGeneric struct {
	Task
	Result chan error
}

func (*TaskGeneric) task() {}

func NewTaskSingle() TaskGeneric {
	return TaskGeneric{Result: make(chan error, 1)}
}

type TaskOwn struct {
	TaskGeneric
	Client   client.Client
	Listener listener.Listener
}

type TaskDelete struct {
	TaskGeneric
	Client   client.Client
	Listener listener.Listener
}

type TaskStart struct {
	TaskGeneric
	Client client.Client
}

type TaskStop struct {
	TaskGeneric
	Client client.Client
}
