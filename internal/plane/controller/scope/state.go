package scope

import (
	"context"
	"mals/internal/scope"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

// rw?
type Resource[T any] struct {
	resource   T
	registered bool
	started    bool
}

type Scope struct {
	scope scope.Scope

	children []*Scope

	models xsync.Map[string, *Resource[*config.Model]]
	lsps   xsync.Map[string, *Resource[*config.Lsp]]
}

type State struct {
	statusRW     sync.RWMutex
	statusCancel context.CancelFunc

	configModels *xsync.Map[string, *config.Model]
	configLsps   *xsync.Map[string, *config.Lsp]
	configUsages *xsync.Map[string, *config.Usage]

	scope *Scope
}
