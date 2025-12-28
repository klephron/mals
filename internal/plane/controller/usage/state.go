package usage

import (
	"context"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type State struct {
	statusRW     sync.RWMutex
	statusCancel context.CancelFunc

	configModels *xsync.Map[string, *config.Model]
	configLsps   *xsync.Map[string, *config.Lsp]
	configUsages *xsync.Map[string, *config.Usage]
}
