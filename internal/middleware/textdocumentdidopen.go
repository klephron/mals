package middleware

import "mals/internal/lsp/protocol"

func (s *Middleware) TextDocumentDidOpen(params *protocol.DidOpenTextDocumentParams) error {
	uri := params.TextDocument.URI

	workspaces := s.workspaceFindAllByPrefix(uri)

	if len(workspaces) == 0 {
		s.plane.Warnf("%v: file %v is not bound to any workspace", s.Name(), uri)
	}

	for _, workspace := range workspaces {
		s.documentAdd(workspace, params.TextDocument.URI, &params.TextDocument.Text)
	}

	return nil
}
