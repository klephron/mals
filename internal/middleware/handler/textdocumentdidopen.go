package handler

import (
	"mals/pkg/config"
	"mals/third_party/lsp"
)

func (s *Handler) TextDocumentDidOpenDefault(params *lsp.DidOpenTextDocumentParams) error {
	uri := params.TextDocument.URI
	workspaces := s.workspaceFindAllByPrefix(uri)

	if len(workspaces) == 0 {
		s.plane.Warnf("%T %v: file %v is not bound to any workspace", s, s.Name(), uri)
	}
	for _, workspace := range workspaces {
		workspace.DocumentAdd(params.TextDocument.URI, &params.TextDocument.Text, params.TextDocument.Version)
	}

	s.resources.Range(func(key string, value *config.HandlerLspResource) bool {
		scope, err := s.getResourceScope(value.Scope)
		if err != nil {
			s.plane.Errorf("%T %v: TextDocumentDidOpen %v", s, s.Name(), err)
			return true
		}

		switch vs := value.Spec.(type) {
		case *config.HandlerLspResourceSpecLsp:
			lspName, token, err := s.plane.Scope().LspAcquire(vs.Name, scope)
			if err != nil {
				s.plane.Errorf("%T %v: TextDocumentDidOpen %v", s, s.Name(), err)
				return true
			}

			defer func() {
				if err := s.plane.Scope().LspRelease(lspName, token); err != nil {
					s.plane.Errorf("%T %v: TextDocumentDidOpen %v", s, s.Name(), err)
				}
			}()

			lspParams := &lsp.DidOpenTextDocumentParams{
				TextDocument: lsp.TextDocumentItem{
					URI:        params.TextDocument.URI,
					LanguageID: params.TextDocument.LanguageID,
					Version:    params.TextDocument.Version,
					Text:       params.TextDocument.Text,
				},
			}

			err = s.plane.Lsp().TextDocumentDidOpen(lspName, lspParams)
			if err != nil {
				s.plane.Errorf("%T %v: TextDocumentDidOpen %v", s, s.Name(), err)
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

func (s *Handler) TextDocumentDidOpen(params *lsp.DidOpenTextDocumentParams) error {
	if *s.endpoints.TextDocumentDidOpen.Default {
		err := s.TextDocumentDidOpenDefault(params)
		if err != nil {
			return err
		}
	}

	s.plane.Infof("%T %v: TextDocumentDidOpen done", s, s.Name())

	return nil
}
