package controller

type LifecycleController interface {
	Serve() error

	Shutdown()
}
