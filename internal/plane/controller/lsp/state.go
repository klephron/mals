package lsp

import (
	"context"
	"mals/internal/lsp/server"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type State struct {
	statusRW     sync.RWMutex
	statusCancel context.CancelFunc

	lsps *xsync.Map[string, *Lsp]
}

type Lsp struct {
	rw         sync.RWMutex
	config     *config.Lsp
	lsp        server.LspServer
	cancelFunc context.CancelFunc
}
