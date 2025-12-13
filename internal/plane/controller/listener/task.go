package listener

import (
	"mals/internal/client"
	"mals/pkg/config"
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

type TaskRegister struct {
	TaskGeneric
	config config.Listener
}

type TaskUnregister struct {
	TaskGeneric
	name string
}

type TaskCreate struct {
	TaskGeneric
	name string
}

type TaskDelete struct {
	TaskGeneric
	name string
}

type TaskStart struct {
	TaskGeneric
	name string
}

type TaskStop struct {
	TaskGeneric
	name string
}

type TaskClientAdd struct {
	TaskGeneric
	name   string
	client client.Client
}

type TaskClientRemove struct {
	TaskGeneric
	name   string
	client client.Client
}
