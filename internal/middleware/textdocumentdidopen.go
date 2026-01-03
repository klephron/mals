package middleware

import (
	"fmt"
	"mals/internal/lsp/protocol"
	"mals/internal/scope"
	"mals/internal/usage"
	"mals/internal/util"
	"mals/pkg/config"
)

func (s *Middleware) eventTextDocumentDidOpenLsp(params *protocol.DidOpenTextDocumentParams, _ *Workspace, step *config.Step) error {
	if step.Scope != "client" {
		s.plane.Warnf("TextDocumentDidOpen %T %v scope %v unsupported, set to client", step, step, step.Scope)
	}
	scope := scope.NewScopeClient(s.client.Name())

	lspName := step.Kind.(*config.StepKindLsp).Name

	lspKey, token, err := s.plane.ScopeLspAcquire(lspName, scope)
	if err != nil {
		s.plane.Errorf("TextDocumentDidOpen %T %v: %v", step, step, err)
		return err
	}
	defer s.plane.ScopeLspRelease(lspKey, token)

	lspParams := &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        params.TextDocument.URI,
			LanguageID: params.TextDocument.LanguageID,
			Version:    params.TextDocument.Version,
			Text:       params.TextDocument.Text,
		},
	}

	err = s.plane.LspEventTextDocumentDidOpen(lspKey, lspParams)
	if err != nil {
		s.plane.Errorf("TextDocumentDidOpen %T %v: %v", step, step, err)
		return nil
	}

	s.plane.Infof("TextDocumentDidOpen %T %v", step, step)

	return nil
}

func (s *Middleware) eventTextDocumentDidOpenWorkflow(params *protocol.DidOpenTextDocumentParams, workspace *Workspace, workflow *config.Workflow) error {
	for _, step := range workflow.Steps {
		switch step.Kind.(type) {
		case *config.StepKindLsp:
			if err := s.eventTextDocumentDidOpenLsp(params, workspace, step); err != nil {
				return err
			}
		default:
			err := fmt.Errorf("TextDocumentDidOpen unhandled %T %v", step, step)
			s.plane.Warnf("%v", err)
			return err
		}
	}

	return nil
}

func (s *Middleware) eventTextDocumentDidOpen(params *protocol.DidOpenTextDocumentParams, workspaces []*Workspace) error {
	for _, workspace := range workspaces {
		usages := s.plane.UsageGetFilteredClient(
			usage.ConditionFilter{Filetype: nil, Path: &workspace.uri},
			usage.EventFilter{Event: util.Ptr(config.EventTextDocumentDidOpen)}, s.client.Name())

		for _, usage := range usages {
			if err := s.eventTextDocumentDidOpenWorkflow(params, workspace, usage.Workflow); err != nil {
				continue
			}
			s.plane.Infof("TextDocumentDidOpen %T %v: usage %v ok", s, s.Name(), usage.Name)
		}
	}
	return nil
}

func (s *Middleware) TextDocumentDidOpen(params *protocol.DidOpenTextDocumentParams) error {
	uri := params.TextDocument.URI

	workspaces := s.workspaceFindAllByPrefix(uri)

	if len(workspaces) == 0 {
		s.plane.Warnf("%v: file %v is not bound to any workspace", s.Name(), uri)
	}

	for _, workspace := range workspaces {
		s.documentAdd(workspace, params.TextDocument.URI, &params.TextDocument.Text, params.TextDocument.Version)
	}

	s.eventTextDocumentDidOpen(params, workspaces)

	return nil
}
