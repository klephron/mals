package model

import (
	"context"
	"mals/internal/model"
	"mals/internal/plane/event"
	"mals/pkg/config"

	"github.com/puzpuzpuz/xsync/v4"
)

type State struct {
	models    *xsync.Map[string, *ModelValue]
	eventChan <-chan event.Event
	taskChan  chan Task
}

type ModelValue struct {
	config     config.Model
	model      model.Model
	cancelFunc context.CancelFunc
	queue      xsync.UMPSCQueue[Task]
}
