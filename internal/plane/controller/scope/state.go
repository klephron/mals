package scope

import (
	"context"
	"mals/internal/scope"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type ResourceLsp struct {
	rw           sync.RWMutex
	fullname     string
	config       *config.Lsp
	dependencies *xsync.Map[string, *ScopeToken]
}

type ResourceModel struct {
	rw           sync.Mutex
	fullname     string
	config       *config.Model
	dependencies *xsync.Map[string, *ScopeToken]
}

type Space struct {
	space scope.Space

	children *xsync.Map[scope.Space, *Space]

	lsps   *xsync.Map[string, *ResourceLsp]
	models *xsync.Map[string, *ResourceModel]
}

func newSpace(space scope.Space) *Space {
	return &Space{
		space:    space,
		children: xsync.NewMap[scope.Space, *Space](),
		lsps:     xsync.NewMap[string, *ResourceLsp](),
		models:   xsync.NewMap[string, *ResourceModel](),
	}
}

func newSpaceRoot() *Space {
	return newSpace(scope.NewSpace(""))
}

type State struct {
	statusRW     sync.RWMutex
	statusCancel context.CancelFunc

	lsps   *xsync.Map[string, *config.Lsp]
	models *xsync.Map[string, *config.Model]

	root *Space
}
