package middleware

import (
	"mals/internal/client"
	"mals/internal/info"
	"mals/internal/lsp/protocol"
)

func (s *Middleware) Initialize(params *protocol.InitializeParams, client client.Client) (*protocol.InitializeResult, error) {

	s.client = client

	for _, workspace := range params.WorkspaceFolders {
		s.workspaceAdd(workspace.URI, workspace.Name)
	}

	s.plane.Infof("%v: initialized", s.Name())

	result := &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    protocol.Full,
			},
			CompletionProvider: &protocol.CompletionOptions{},
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    info.MiddlewareServerName,
			Version: info.MiddlewareVersion,
		},
	}

	return result, nil
}
