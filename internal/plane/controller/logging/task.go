package logging

import (
	"mals/internal/log"
	"mals/pkg/config"
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

type TaskShutdown struct {
	TaskGeneric
}

type TaskTerminate struct {
	TaskGeneric
}

type TaskRegister struct {
	TaskGeneric
	Config config.Log
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

type TaskLog struct {
	TaskGeneric
	Level log.Level
	Msg   string
}
