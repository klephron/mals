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

	usages *xsync.Map[string, *config.Usage]
}
