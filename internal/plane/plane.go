package plane

import "mals/internal/plane/controller"

type Plane interface {
	Client() controller.ClientController
	Listener() controller.ListenerController
	Log() controller.LogController

	Serve(onReady func())
	Shutdown() error
}
