package middleware

import (
	"mals/internal/info"
	"mals/internal/lsp/protocol"
)

func (s *Middleware) Initialize(params *protocol.InitializeParams) (*protocol.ServerCapabilities, *protocol.ServerInfo, error) {
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

	return capabilities, serverInfo, nil
}
