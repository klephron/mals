package middleware

import "mals/internal/lsp/protocol"

func (s *Middleware) TextDocumentDidChange(params *protocol.DidChangeTextDocumentParams) error {
	uri := params.TextDocument.URI

	workspaces := s.workspaceFindAll(params.TextDocument.URI)

	if len(workspaces) == 0 {
		s.plane.Warnf("%v: file %v is not bound to any workspace", s.Name(), uri)
	}

	for _, workspace := range workspaces {
		document := s.documentGet(workspace, params.TextDocument.URI)
		if document == nil {
			continue
		}

		var text *string
		for _, change := range params.ContentChanges {
			if change.Range != nil {
				s.plane.Warnf("%s: unsupported ContentChange Range is unsupported", s.Name())
				continue
			}

			text = &change.Text
		}

		if text != nil {
			s.documentUpdateFull(workspace, document, text, params.TextDocument.Version)
		}
	}

	return nil
}
