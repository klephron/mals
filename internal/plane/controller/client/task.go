package client

import (
	"mals/internal/client"
)

type Task interface {
	task()
}

type TaskGeneric struct {
	Task
	result chan error
}

func (*TaskGeneric) task() {}

func newTask() TaskGeneric {
	return TaskGeneric{result: make(chan error, 1)}
}

type TaskShutdown struct {
	TaskGeneric
}

type TaskOwn struct {
	TaskGeneric
	client   client.Client
	listener string
}

type TaskDelete struct {
	TaskGeneric
	client client.Client
	notify bool
}

type TaskStart struct {
	TaskGeneric
	client client.Client
}

type TaskStop struct {
	TaskGeneric
	client client.Client
}
