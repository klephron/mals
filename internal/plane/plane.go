package plane

import (
	"mals/internal/plane/controller"
)

type Plane interface {
	Run(onReady func())
	Shutdown() error

	Listener() controller.ListenerController
	Log() controller.LogController
	Lsp() controller.LspController
	Model() controller.ModelController
	Scope() controller.ScopeController
	Usage() controller.UsageController

	Debugf(format string, a ...any) error
	Infof(format string, a ...any) error
	Warnf(format string, a ...any) error
	Errorf(format string, a ...any) error
}
