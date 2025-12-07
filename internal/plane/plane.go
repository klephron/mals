package plane

import "mals/internal/plane/controller"

type Plane interface {
	Serve(onReady func())
	Shutdown() error
	Terminate() error

	Client() controller.ClientController
	Listener() controller.ListenerController
	Log() controller.LogController
}
