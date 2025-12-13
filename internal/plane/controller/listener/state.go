package listener

import (
	"context"
	"mals/internal/client"
	"mals/internal/listener"
	"mals/internal/plane/event"
	"mals/pkg/config"

	"github.com/puzpuzpuz/xsync/v4"
)

type State struct {
	listeners *xsync.Map[string, *ListenerValue]
	external  <-chan event.Event
	internal  chan Task
}

type ListenerValue struct {
	config     config.Listener
	listener   listener.Listener
	cancelFunc context.CancelFunc
	clients    *xsync.Map[client.Client, struct{}]
}
