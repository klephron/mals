package listener

import (
	"context"
	"mals/internal/listener"
	"mals/internal/listener/api"
	"mals/internal/listener/lsp"
	"mals/pkg/config"
	"sync"

	"github.com/puzpuzpuz/xsync/v4"
)

type State struct {
	statusRW     sync.RWMutex
	statusCancel context.CancelFunc
	listeners    *xsync.Map[string, *Listener]
}

type Listener struct {
	rw         sync.RWMutex
	config     *config.Listener
	cancelFunc context.CancelFunc
	listener   ListenerMixin
}

type ListenerMixin interface {
	Listener() listener.Listener
}

type ListenerMixinApi struct {
	listener *api.ListenerApi
}

func (s *ListenerMixinApi) Listener() listener.Listener {
	return s.listener
}

type ListenerMixinLsp struct {
	listener *lsp.ListenerLsp
	clients  *xsync.Map[string, *ListenerLspClient]
}

func (s *ListenerMixinLsp) Listener() listener.Listener {
	return s.listener
}

type ListenerLspClient struct {
	client     listener.ListenerLspClient
	cancelFunc context.CancelFunc
}
