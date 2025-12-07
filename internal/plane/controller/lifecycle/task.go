package lifecycle

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
