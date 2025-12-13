package client

import (
	"context"
	"mals/internal/client"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type ClientValue struct {
	rw         sync.RWMutex
	listener   string
	cancelFunc context.CancelFunc
}

type State struct {
	statusCancel context.CancelFunc
	statusRW     sync.RWMutex

	clients *xsync.Map[client.Client, *ClientValue]
}
