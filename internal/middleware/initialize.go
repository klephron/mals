package middleware

import (
	"fmt"
	"mals/internal/client"
	"mals/internal/info"
	"mals/internal/lsp/protocol"
	"mals/internal/scope"
	"mals/internal/usage"
	"mals/internal/util"
	"mals/pkg/config"
	"os"
)

func (s *Middleware) eventInitializeLsp(params *protocol.InitializeParams, workspace *Workspace, step *config.Step) error {
	if step.Scope != "client" {
		s.plane.Warnf("Initialize step %T %v scope %v unsupported", step, step, step.Scope)
	}

	lspName := step.Kind.(*config.StepKindLsp).Name

	scope := scope.NewScopeClient(s.client.Name())
	lspKey, token, err := s.plane.ScopeLspAcquire(lspName, scope)

	if err != nil {
		s.plane.Warnf("Initialize step %T %v: %v", step, step, err)
		return err
	}
	defer s.plane.ScopeLspRelease(lspKey, token)

	lspParams := &protocol.InitializeParams{
		XInitializeParams: protocol.XInitializeParams{
			ProcessID: int32(os.Getpid()),
			ClientInfo: &protocol.ClientInfo{
				Name:    info.MiddlewareClientName,
				Version: info.MiddlewareVersion,
			},
			Locale:                params.Locale,
			Capabilities:          params.Capabilities,
			InitializationOptions: params.InitializationOptions,
			Trace:                 params.Trace,
		},
		WorkspaceFoldersInitializeParams: protocol.WorkspaceFoldersInitializeParams{
			WorkspaceFolders: []protocol.WorkspaceFolder{
				{
					URI:  workspace.uri,
					Name: workspace.name,
				},
			},
		},
	}

	result, err := s.plane.LspEventInitialize(lspKey, lspParams)
	if err != nil {
		s.plane.Warnf("Initialize step %T %v: %v", step, step, err)
		return nil
	}

	s.plane.Infof("Initialize step %T %v: %+v", step, step, result)

	return nil
}

func (s *Middleware) eventInitializeWorkflow(params *protocol.InitializeParams, workspace *Workspace, workflow *config.Workflow) error {
	for _, step := range workflow.Steps {
		switch step.Kind.(type) {
		case *config.StepKindLsp:
			if err := s.eventInitializeLsp(params, workspace, step); err != nil {
				return err
			}
		default:
			err := fmt.Errorf("Initialize unhandled step %T %v", step, step)
			return err
		}
	}

	return nil
}

func (s *Middleware) eventInitialize(params *protocol.InitializeParams, workspaces []*Workspace) error {
	for _, workspace := range workspaces {
		usages := s.plane.UsageGetFilteredClient(
			usage.ConditionFilter{Filetype: nil, Path: &workspace.uri},
			usage.EventFilter{Event: util.Ptr(config.EventInitialize)}, s.client.Name())

		for _, usage := range usages {
			s.plane.Infof("%T %v: usage %v: Initialize", s, s.Name(), usage.Name)

			err := s.eventInitializeWorkflow(params, workspace, usage.Workflow)
			if err != nil {
				s.plane.Warnf("%v", err)
			}
		}
	}

	return nil
}

func (s *Middleware) Initialize(params *protocol.InitializeParams, client client.Client) (*protocol.InitializeResult, error) {

	s.client = client

	for _, workspace := range params.WorkspaceFolders {
		s.workspaceAdd(workspace.URI, workspace.Name)
	}

	workspaces := s.workspaceFindAll()

	s.eventInitialize(params, workspaces)

	s.plane.Infof("%v: Initialized", s.Name())

	result := &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    protocol.Full,
			},
			CompletionProvider: &protocol.CompletionOptions{},
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    info.MiddlewareServerName,
			Version: info.MiddlewareVersion,
		},
	}

	return result, nil
}
