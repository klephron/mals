package model

import (
	"context"
	"mals/internal/model/queued"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type State struct {
	statusRW     sync.RWMutex
	statusCancel context.CancelFunc

	models *xsync.Map[string, *Model]
}

type Model struct {
	rw         sync.RWMutex
	config     *config.Model
	model      *queued.ModelQueued
	cancelFunc context.CancelFunc
}
