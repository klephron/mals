package middleware

import "mals/internal/lsp/protocol"

func (s *Middleware) TextDocumentDidChange(params *protocol.DidChangeTextDocumentParams) error {
	uri := params.TextDocument.URI

	workspaces := s.workspaceFindAllByPrefix(params.TextDocument.URI)

	if len(workspaces) == 0 {
		s.plane.Warnf("%v: file %v is not bound to any workspace", s.Name(), uri)
	}

	for _, workspace := range workspaces {
		document := s.documentGet(workspace, params.TextDocument.URI)
		if document == nil {
			continue
		}

		s.documentUpdate(workspace, document, params.TextDocument.Version, params.ContentChanges)
	}

	return nil
}
