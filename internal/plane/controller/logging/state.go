package logging

import (
	"context"
	"mals/internal/log"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type LogValue struct {
	rw      sync.RWMutex
	config  *config.Log
	log     log.Log
	enabled bool
}

type State struct {
	statusCancel context.CancelFunc
	statusRW     sync.RWMutex

	logs *xsync.Map[string, *LogValue]
}
