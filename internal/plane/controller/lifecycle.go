package controller

type LifecycleController interface {
	Serve(onReady func()) error

	Shutdown() error
}
