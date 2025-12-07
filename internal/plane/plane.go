package plane

import "mals/internal/plane/controller"

type Plane interface {
	Client() controller.ClientController
	Lifecycle() controller.LifecycleController
	Listener() controller.ListenerController
	Log() controller.LogController

	Serve(onReady func())
}
