package middleware

import "mals/internal/lsp/protocol"

func (s *Middleware) TextDocumentDidClose(params *protocol.DidCloseTextDocumentParams) error {
	uri := params.TextDocument.URI

	workspaces := s.workspaceFindAllByPrefix(uri)

	if len(workspaces) == 0 {
		s.plane.Warnf("%v: file %v is not bound to any workspace", s.Name(), uri)
	}

	for _, workspace := range workspaces {
		document := s.documentGet(workspace, params.TextDocument.URI)
		if document == nil {
			continue
		}

		s.documentDelete(workspace, document.Uri())
	}

	return nil
}
