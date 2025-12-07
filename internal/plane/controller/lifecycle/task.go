package lifecycle

type Task interface {
	task()
}

type TaskGeneric struct {
	Task
}

func (*TaskGeneric) task() {}

type TaskShutdown struct {
	TaskGeneric
}

type TaskTerminate struct {
	TaskGeneric
}
