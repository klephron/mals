package middleware

import (
	"fmt"
	"mals/internal/lsp/protocol"
	"mals/internal/scope"
	"mals/internal/usage"
	"mals/internal/util"
	"mals/pkg/config"
)

func (s *Middleware) eventTextDocumentDidChangeLsp(params *protocol.DidChangeTextDocumentParams, workspace *Workspace, step *config.Step) error {
	if step.Scope != "client" {
		s.plane.Warnf("TextDocumentDidChange %T %v scope %v unsupported, set to client", step, step, step.Scope)
	}
	scope := scope.NewScopeClient(s.client.Name())

	lspName := step.Kind.(*config.StepKindLsp).Name

	lspKey, token, err := s.plane.ScopeLspAcquire(lspName, scope)
	if err != nil {
		s.plane.Errorf("TextDocumentDidChange %T %v: %v", step, step, err)
		return err
	}
	defer s.plane.ScopeLspRelease(lspKey, token)

	capabilities, err := s.plane.LspGetCapabilities(lspKey)
	if err != nil {
		s.plane.Errorf("TextDocumentDidChange %T %v: %v", step, step, err)
		return err
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
			s.plane.Errorf("TextDocumentDidChange %T %v: %v", step, step, err)
			return err
		}
		if syncOptions, err := util.JsonUnmarshal[protocol.TextDocumentSyncOptions](data); err != nil {
			s.plane.Warnf("TextDocumentDidChange %T %v: %v", step, step, err)
			syncKind = syncOptions.Change
		} else {
			syncKind = syncOptions.Change
		}
	default:
		err := fmt.Errorf("%v: capabilities.TextDocumentSync has unexpected type %T", lspKey, v)
		s.plane.Errorf("TextDocumentDidChange %T %v: %v", step, step, err)
		return err
	}

	s.plane.Debugf("TextDocumentDidChange %T %v: sync kind %d", step, step, syncKind)

	// works for middleware in incremental mode:
	// s.textDocumentSyncKind == protocol.Incremental

	var lspParams *protocol.DidChangeTextDocumentParams

	switch syncKind {
	case protocol.Full:
		uri := params.TextDocument.TextDocumentIdentifier.URI

		document := s.documentGet(workspace, uri)
		if document == nil {
			return nil
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
		err := fmt.Errorf("%v: unhandled sync kind %d", lspKey, syncKind)
		s.plane.Errorf("TextDocumentDidChange %T %v: %v", step, step, err)
		return err
	}

	err = s.plane.LspEventTextDocumentDidChange(lspKey, lspParams)
	if err != nil {
		s.plane.Errorf("TextDocumentDidChange %T %v: %v", step, step, err)
		return nil
	}

	s.plane.Debugf("TextDocumentDidChange %T %v", step, step)

	return nil
}

func (s *Middleware) eventTextDocumentDidChangeWorkflow(params *protocol.DidChangeTextDocumentParams, workspace *Workspace, workflow *config.Workflow) error {
	for _, step := range workflow.Steps {
		switch step.Kind.(type) {
		case *config.StepKindLsp:
			if err := s.eventTextDocumentDidChangeLsp(params, workspace, step); err != nil {
				return err
			}
		default:
			err := fmt.Errorf("TextDocumentDidChange unhandled %T %v", step, step)
			s.plane.Warnf("%v", err)
			return err
		}
	}

	return nil
}

func (s *Middleware) eventTextDocumentDidChange(params *protocol.DidChangeTextDocumentParams, workspaces []*Workspace) error {
	for _, workspace := range workspaces {
		usages := s.plane.UsageGetFilteredClient(
			usage.ConditionFilter{Filetype: nil, Path: &workspace.uri},
			usage.EventFilter{Event: util.Ptr(config.EventTextDocumentDidOpen)}, s.client.Name())

		for _, usage := range usages {
			if err := s.eventTextDocumentDidChangeWorkflow(params, workspace, usage.Workflow); err != nil {
				continue
			}
			s.plane.Infof("TextDocumentDidChange %T %v: usage %v ok", s, s.Name(), usage.Name)
		}
	}
	return nil
}

func (s *Middleware) TextDocumentDidChange(params *protocol.DidChangeTextDocumentParams) error {
	uri := params.TextDocument.URI
	workspaces := s.workspaceFindAllByPrefix(uri)

	if len(workspaces) == 0 {
		s.plane.Warnf("%v: file %v is not bound to any workspace", s.Name(), uri)
	}

	for _, workspace := range workspaces {
		document := s.documentGet(workspace, uri)
		if document == nil {
			continue
		}

		s.documentUpdate(workspace, document, params.TextDocument.Version, params.ContentChanges)
	}

	s.eventTextDocumentDidChange(params, workspaces)

	return nil
}
