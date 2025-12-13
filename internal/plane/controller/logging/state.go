package logging

import (
	"mals/internal/log"
	"mals/internal/plane/event"
	"mals/pkg/config"

	"github.com/puzpuzpuz/xsync/v4"
)

type State struct {
	logs      *xsync.Map[string, *LogValue]
	eventChan <-chan event.Event
	taskChan  chan Task
}

type LogValue struct {
	config  config.Log
	log     log.Log
	enabled bool
}
