package client

import (
	"context"
	"mals/internal/client"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type ClientValue struct {
	rw         sync.RWMutex
	client     client.Client
	listener   string
	cancelFunc context.CancelFunc
}

type State struct {
	statusRW     sync.RWMutex
	statusCancel context.CancelFunc

	clients *xsync.Map[string, *ClientValue]
}
