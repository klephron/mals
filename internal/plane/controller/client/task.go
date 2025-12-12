package client

import (
	"mals/internal/client"
)

type Task interface {
	task()
}

type TaskGeneric struct {
	Task
	Result chan error
}

func (*TaskGeneric) task() {}

func NewTask() TaskGeneric {
	return TaskGeneric{Result: make(chan error, 1)}
}

type TaskShutdown struct {
	TaskGeneric
}

type TaskOwn struct {
	TaskGeneric
	Client   client.Client
	Listener string
}

type TaskDelete struct {
	TaskGeneric
	Client client.Client
	Notify bool
}

type TaskStart struct {
	TaskGeneric
	Client client.Client
}

type TaskStop struct {
	TaskGeneric
	Client client.Client
}
