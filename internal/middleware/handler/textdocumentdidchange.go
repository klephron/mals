package handler

import (
	"fmt"
	"mals/internal/lsp/protocol"
	"mals/internal/middleware/document"
	"mals/internal/util"
	"mals/pkg/config"
)

func (s *Handler) TextDocumentDidChangeDefault(params *protocol.DidChangeTextDocumentParams) error {
	uri := params.TextDocument.URI
	workspaces := s.workspaceFindAllByPrefix(uri)

	if len(workspaces) == 0 {
		s.plane.Warnf("%T %v: TextDocumentDidChange file %v not bound to workspace", s, s.Name(), uri)
	}

	for _, workspace := range workspaces {
		document := s.documentGet(workspace, uri)
		if document == nil {
			continue
		}

		s.documentUpdate(workspace, document, params.TextDocument.Version, params.ContentChanges)
	}

	s.resources.Range(func(key string, value *config.HandlerLspResource) bool {
		scope, err := s.getResourceScope(value.Scope)
		if err != nil {
			s.plane.Errorf("%T %v: TextDocumentDidChange %v", s, s.Name(), err)
			return true
		}

		switch vs := value.Spec.(type) {
		case *config.HandlerLspResourceSpecLsp:
			lspName, token, err := s.plane.Scope().LspAcquire(vs.Name, scope)
			if err != nil {
				s.plane.Errorf("%T %v: TextDocumentDidChange %v", s, s.Name(), err)
				return true
			}

			defer func() {
				if err := s.plane.Scope().LspRelease(lspName, token); err != nil {
					s.plane.Errorf("%T %v: TextDocumentDidChange %v", s, s.Name(), err)
				}
			}()

			capabilities, err := s.plane.Lsp().GetCapabilities(lspName)
			if err != nil {
				s.plane.Errorf("%T %v: TextDocumentDidChange %v", s, s.Name(), err)
				return true
			}

			var syncKind protocol.TextDocumentSyncKind

			switch v := capabilities.TextDocumentSync.(type) {
			case protocol.TextDocumentSyncOptions:
				syncKind = v.Change
			case protocol.TextDocumentSyncKind:
				syncKind = v
			case float64:
				syncKind = protocol.TextDocumentSyncKind(v)
			case map[string]any:
				data, err := util.JsonMarshal(&v)
				if err != nil {
					s.plane.Errorf("%T %v: TextDocumentDidChange %v", s, s.Name(), err)
					return true
				}
				if syncOptions, err := util.JsonUnmarshal[protocol.TextDocumentSyncOptions](data); err != nil {
					s.plane.Warnf("%T %v: TextDocumentDidChange %v", s, s.Name(), err)
					syncKind = syncOptions.Change
				} else {
					syncKind = syncOptions.Change
				}
			default:
				s.plane.Errorf("%T %v: TextDocumentDidChange capabilities.TextDocumentSync has unexpected type %T", s, s.Name(), v)
				return true
			}

			s.plane.Debugf("%T %v: TextDocumentDidChange sync kind %d", s, s.Name(), syncKind)

			var lspParams *protocol.DidChangeTextDocumentParams

			switch syncKind {
			case protocol.Full:
				uri := params.TextDocument.TextDocumentIdentifier.URI

				var document *document.Document

				for _, workspace := range s.workspaceFindAllByPrefix(params.TextDocument.URI) {
					document = s.documentGet(workspace, uri)
					if document != nil {
						break
					}
				}

				if document == nil {
					s.plane.Errorf("%T %v: TextDocumentDidChange document %v not found", uri)
					return true
				}

				lspParams = &protocol.DidChangeTextDocumentParams{
					TextDocument: protocol.VersionedTextDocumentIdentifier{
						TextDocumentIdentifier: protocol.TextDocumentIdentifier{
							URI: uri,
						},
						Version: params.TextDocument.Version,
					},
					ContentChanges: []protocol.TextDocumentContentChangeEvent{
						{
							Text: document.Text(),
						},
					},
				}

			case protocol.Incremental:
				lspParams = params

			default:
				err := fmt.Errorf("%v: unhandled sync kind %d", lspName, syncKind)
				s.plane.Errorf("%T %v: TextDocumentDidChange %v", s, s.Name(), err)
				return true
			}

			err = s.plane.Lsp().TextDocumentDidChange(lspName, lspParams)
			if err != nil {
				s.plane.Errorf("%T %v: TextDocumentDidChange %v", s, s.Name(), err)
				return true
			}

			s.plane.Debugf("%T %v: TextDocumentDidChange", s, s.Name())

		case *config.HandlerLspResourceSpecModel:
		default:
			s.plane.Errorf("%T %v: TextDocumentDidChange unexpected spec %T", s, s.Name(), vs)
		}

		return true
	})

	return nil
}

func (s *Handler) TextDocumentDidChange(params *protocol.DidChangeTextDocumentParams) error {
	if *s.endpoints.TextDocumentDidChange.Default {
		err := s.TextDocumentDidChangeDefault(params)
		if err != nil {
			return err
		}
	}

	s.plane.Infof("%T %v: TextDocumentDidChange done", s, s.Name())

	return nil
}
