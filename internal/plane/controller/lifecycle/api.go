package lifecycle

func (s *LifecycleController) Shutdown() {
	s.internal <- &EventShutdown{}
}

func (s *LifecycleController) Terminate() {
	s.internal <- &EventTerminate{}
}
