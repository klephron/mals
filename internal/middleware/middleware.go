package middleware

import (
	"fmt"
	"mals/internal/lsp/protocol"
	"mals/internal/middleware/document"
	"mals/internal/plane"

	"github.com/puzpuzpuz/xsync/v4"
)

type workspace struct {
	uri       string
	name      string
	documents *xsync.Map[string, *document.Document]
}

type Middleware struct {
	plane plane.Plane

	listenerName string
	clientName   string

	initialized bool
	workspaces  *xsync.Map[string, *workspace]

	textDocumentSyncKind protocol.TextDocumentSyncKind
}

func New(plane plane.Plane) *Middleware {
	return &Middleware{
		plane:                plane,
		listenerName:         "",
		clientName:           "",
		initialized:          false,
		workspaces:           xsync.NewMap[string, *workspace](),
		textDocumentSyncKind: protocol.Incremental,
	}
}

func (s *Middleware) Name() string {
	return fmt.Sprintf("%v:%v", s.listenerName, s.clientName)
}
