package middleware

import (
	"mals/internal/client"
	"mals/internal/info"
	"mals/internal/lsp/protocol"
)

func (s *Middleware) Initialize(params *protocol.InitializeParams, client client.Client) (*protocol.ServerCapabilities, *protocol.ServerInfo, error) {

	s.client = client

	for _, workspace := range params.WorkspaceFolders {
		s.workspaceAdd(workspace.URI, workspace.Name)
	}

	capabilities := &protocol.ServerCapabilities{
		TextDocumentSync: protocol.TextDocumentSyncOptions{
			OpenClose: true,
			Change:    protocol.Full,
		},
		CompletionProvider: &protocol.CompletionOptions{},
	}

	serverInfo := &protocol.ServerInfo{
		Name:    info.MiddlewareServerName,
		Version: info.MiddlewareVersion,
	}

	s.plane.Infof("%v: initialized", s.Name())

	return capabilities, serverInfo, nil
}
