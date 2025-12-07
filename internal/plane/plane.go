package plane

import "mals/internal/plane/controller"

type Plane interface {
	Lifecycle() controller.LifecycleController
	Listener() controller.ListenerController
	Log() controller.LogController

	Serve(onReady func())
}
