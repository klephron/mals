package handler

import (
	"mals/internal/lsp/protocol"
	"mals/pkg/config"
)

func (s *Handler) TextDocumentDidCloseDefault(params *protocol.DidCloseTextDocumentParams) error {
	uri := params.TextDocument.URI
	workspaces := s.workspaceFindAllByPrefix(uri)

	if len(workspaces) == 0 {
		s.plane.Warnf("%T %v: TextDocumentDidClose file %v not bound to workspace", s, s.Name(), uri)
	}

	for _, workspace := range workspaces {
		document := s.documentGet(workspace, params.TextDocument.URI)
		if document == nil {
			continue
		}

		s.documentDelete(workspace, params.TextDocument.URI)
	}

	s.resources.Range(func(key string, value *config.HandlerLspResource) bool {
		scope, err := s.getResourceScope(value.Scope)
		if err != nil {
			s.plane.Errorf("%T %v: TextDocumentDidClose %v", s, s.Name(), err)
			return true
		}

		switch vs := value.Spec.(type) {
		case *config.HandlerLspResourceSpecLsp:
			lspName, token, err := s.plane.Scope().LspAcquire(vs.Name, scope)
			if err != nil {
				s.plane.Errorf("%T %v: TextDocumentDidClose %v", s, s.Name(), err)
				return true
			}

			defer func() {
				if err := s.plane.Scope().LspRelease(lspName, token); err != nil {
					s.plane.Errorf("%T %v: TextDocumentDidClose %v", s, s.Name(), err)
				}
			}()

			lspParams := &protocol.DidCloseTextDocumentParams{
				TextDocument: protocol.TextDocumentIdentifier{
					URI: params.TextDocument.URI,
				},
			}

			err = s.plane.Lsp().TextDocumentDidClose(lspName, lspParams)
			if err != nil {
				s.plane.Errorf("%T %v: TextDocumentDidClose %v", s, s.Name(), err)
				return true
			}

		case *config.HandlerLspResourceSpecModel:
		default:
			s.plane.Errorf("%T %v: TextDocumentDidOpen unexpected spec %T", s, s.Name(), vs)
		}

		return true
	})

	return nil
}

func (s *Handler) TextDocumentDidClose(params *protocol.DidCloseTextDocumentParams) error {
	if *s.endpoints.TextDocumentDidClose.Default {
		err := s.TextDocumentDidCloseDefault(params)
		if err != nil {
			return err
		}
	}

	s.plane.Infof("%T %v: TextDocumentDidClose done", s, s.Name())

	return nil
}
