package model

import "mals/pkg/config"

type Task interface {
	task()
}

type TaskGeneric struct {
	Task
	result chan error
}

func (*TaskGeneric) task() {}

func NewTask() TaskGeneric {
	return TaskGeneric{result: make(chan error, 1)}
}

type TaskShutdown struct {
	TaskGeneric
}

type TaskRegister struct {
	TaskGeneric
	config config.Model
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
