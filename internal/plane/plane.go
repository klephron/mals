package plane

import "mals/internal/plane/controller"

type Plane interface {
	Serve(onReady func())
	Shutdown() error

	Client() controller.ClientController
	Listener() controller.ListenerController
	Log() controller.LogController
	Model() controller.ModelController
	Scope() controller.ScopeController
	Usage() controller.UsageController
}
