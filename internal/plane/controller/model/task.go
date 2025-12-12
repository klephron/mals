package model

import "mals/pkg/config"

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

type TaskRegister struct {
	TaskGeneric
	Config config.Model
}

type TaskUnregister struct {
	TaskGeneric
	Name string
}

type TaskCreate struct {
	TaskGeneric
	Name string
}

type TaskDelete struct {
	TaskGeneric
	Name string
}

type TaskStart struct {
	TaskGeneric
	Name string
}

type TaskStop struct {
	TaskGeneric
	Name string
}
