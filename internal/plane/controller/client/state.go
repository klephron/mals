package client

import (
	"context"
	"mals/internal/client"
	"mals/internal/plane/event"

	"github.com/puzpuzpuz/xsync/v4"
)

type ClientValue struct {
	listener   string
	cancelFunc context.CancelFunc
}

type State struct {
	clients   *xsync.Map[client.Client, *ClientValue]
	eventChan <-chan event.Event
	taskChan  chan Task
}
