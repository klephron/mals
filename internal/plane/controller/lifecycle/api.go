package lifecycle

func (s *LifecycleController) Shutdown() error {
	e := TaskShutdown{TaskGeneric: NewTaskSingle()}
	s.internal <- &e
	return <-e.Result
}

func (s *LifecycleController) Terminate() error {
	e := TaskTerminate{TaskGeneric: NewTaskSingle()}
	s.internal <- &e
	return <-e.Result
}
