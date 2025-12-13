package listener

import (
	"context"
	"mals/internal/client"
	"mals/internal/listener"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type State struct {
	statusRW     sync.RWMutex
	statusCancel context.CancelFunc
	listeners    *xsync.Map[string, *ListenerValue]
}

type ListenerValue struct {
	rw         sync.RWMutex
	config     config.Listener
	listener   listener.Listener
	cancelFunc context.CancelFunc
	clients    *xsync.Map[client.Client, struct{}]
}
