package handler

import (
// "fmt"
// "mals/internal/lsp/protocol"
// "mals/internal/scope"
// "mals/internal/util"
// "mals/pkg/config"
)

// func (s *Middleware) eventTextDocumentDidCloseLsp(params *protocol.DidCloseTextDocumentParams, _ *Workspace, step *config.Step) error {
// 	if step.Scope != "client" {
// 		s.plane.Warnf("%T %v: TextDocumentDidClose %T %v scope %v unsupported, set to client",
// 			s, s.Name(), step, step, step.Scope)
// 	}
// 	scope := scope.NewScopeClient(s.listenerName, s.clientName)

// 	lspName := step.Kind.(*config.StepKindLsp).Name

// 	lspKey, token, err := s.plane.Scope().LspAcquire(lspName, scope)
// 	if err != nil {
// 		s.plane.Errorf("%T %v: TextDocumentDidClose %T %v: %v", s, s.Name(), step, step, err)
// 		return err
// 	}
// 	defer s.plane.Scope().LspRelease(lspKey, token)

// 	lspParams := &protocol.DidCloseTextDocumentParams{
// 		TextDocument: protocol.TextDocumentIdentifier{
// 			URI: params.TextDocument.URI,
// 		},
// 	}

// 	err = s.plane.Lsp().EventTextDocumentDidClose(lspKey, lspParams)
// 	if err != nil {
// 		s.plane.Errorf("%T %v: TextDocumentDidClose %T %v: %v", s, s.Name(), step, step, err)
// 		return nil
// 	}

// 	s.plane.Debugf("%T %v: TextDocumentDidClose %T %v", s, s.Name(), step, step)

// 	return nil
// }

// func (s *Middleware) eventTextDocumentDidCloseWorkflow(params *protocol.DidCloseTextDocumentParams, workspace *Workspace, workflow *config.Workflow) error {
// 	for _, step := range workflow.Steps {
// 		switch step.Kind.(type) {
// 		case *config.StepKindLsp:
// 			if err := s.eventTextDocumentDidCloseLsp(params, workspace, step); err != nil {
// 				return err
// 			}
// 		default:
// 			err := fmt.Errorf("TextDocumentDidOpen unhandled %T %v", step, step)
// 			s.plane.Warnf("%T %v: %v", s, s.Name(), err)
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (s *Middleware) eventTextDocumentDidClose(params *protocol.DidCloseTextDocumentParams, workspaces []*Workspace) error {
// 	for _, workspace := range workspaces {
// 		usages := s.plane.Usage().GetFilteredClient(
// 			usage.ConditionFilter{Filetype: nil, Path: &workspace.uri},
// 			usage.EventFilter{Event: util.Ptr(config.EventTextDocumentDidClose)}, s.listenerName, s.clientName)

// 		for _, usage := range usages {
// 			if err := s.eventTextDocumentDidCloseWorkflow(params, workspace, usage.Workflow); err != nil {
// 				continue
// 			}
// 			s.plane.Infof("%T %v: TextDocumentDidClose usage %v ok", s, s.Name(), usage.Name)
// 		}
// 	}
// 	return nil
// }

// func (s *Middleware) TextDocumentDidClose(params *protocol.DidCloseTextDocumentParams) error {
// 	uri := params.TextDocument.URI

// 	workspaces := s.workspaceFindAllByPrefix(uri)

// 	if len(workspaces) == 0 {
// 		s.plane.Warnf("%T %v: file %v is not bound to any workspace", s, s.Name(), uri)
// 	}

// 	for _, workspace := range workspaces {
// 		document := s.documentGet(workspace, params.TextDocument.URI)
// 		if document == nil {
// 			continue
// 		}

// 		s.documentDelete(workspace, params.TextDocument.URI)
// 	}

// 	// s.eventTextDocumentDidClose(params, workspaces)

// 	return nil
// }
