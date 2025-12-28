package model

import (
	"context"
	"mals/internal/model"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type ModelValue struct {
	rw         sync.RWMutex
	config     *config.Model
	model      model.Model
	queue      *TaskQueue
	cancelFunc context.CancelFunc
}

type State struct {
	statusRW     sync.RWMutex
	statusCancel context.CancelFunc

	models *xsync.Map[string, *ModelValue]
}
