package logging

import (
	"mals/internal/log"
	"mals/internal/plane/event"
	"mals/pkg/config"

	"github.com/puzpuzpuz/xsync/v4"
)

type State struct {
	logs     *xsync.Map[string, *LogValue]
	external <-chan event.Event
	internal chan Task
}

type LogValue struct {
	config  config.Log
	log     log.Log
	enabled bool
}
