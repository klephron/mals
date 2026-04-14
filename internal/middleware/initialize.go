package middleware

import (
	// "fmt"
	"mals/internal/lsp/protocol"
	// "mals/internal/scope"
	// "mals/internal/util"
	// "mals/pkg/config"
	"mals/pkg/info"
	// "os"
)

// func (s *Middleware) eventInitializeLsp(params *protocol.InitializeParams, workspace *Workspace, step *config.Step) error {
// 	if step.Scope != "client" {
// 		s.plane.Warnf("%T %v: Initialize %T %v scope %v unsupported, set to client",
// 			s, s.Name(), step, step, step.Scope)
// 	}
// 	scope := scope.NewScopeClient(s.listenerName, s.clientName)

// 	lspName := step.Kind.(*config.StepKindLsp).Name

// 	lspKey, token, err := s.plane.Scope().LspAcquire(lspName, scope)
// 	if err != nil {
// 		s.plane.Errorf("%T %v: Initialize %T %v: %v", s, s.Name(), step, step, err)
// 		return err
// 	}
// 	defer s.plane.Scope().LspRelease(lspKey, token)

// 	lspParams := &protocol.InitializeParams{
// 		XInitializeParams: protocol.XInitializeParams{
// 			ProcessID: int32(os.Getpid()),
// 			ClientInfo: &protocol.ClientInfo{
// 				Name:    info.MiddlewareClientName,
// 				Version: info.MiddlewareVersion,
// 			},
// 			Locale: params.Locale,
// 			Capabilities: protocol.ClientCapabilities{
// 				TextDocument: protocol.TextDocumentClientCapabilities{
// 					Synchronization: &protocol.TextDocumentSyncClientCapabilities{
// 						DynamicRegistration: false,
// 						WillSave:            false,
// 						WillSaveWaitUntil:   false,
// 						DidSave:             false,
// 					},
// 					Completion: protocol.CompletionClientCapabilities{
// 						DynamicRegistration: false,
// 					},
// 				},
// 			},
// 			InitializationOptions: params.InitializationOptions,
// 			Trace:                 params.Trace,
// 		},
// 		WorkspaceFoldersInitializeParams: protocol.WorkspaceFoldersInitializeParams{
// 			WorkspaceFolders: []protocol.WorkspaceFolder{
// 				{
// 					URI:  workspace.uri,
// 					Name: workspace.name,
// 				},
// 			},
// 		},
// 	}

// 	result, err := s.plane.Lsp().EventInitialize(lspKey, lspParams)
// 	if err != nil {
// 		s.plane.Errorf("%T %v: Initialize %T %v: %v", s, s.Name(), step, step, err)
// 		return nil
// 	}

// 	s.plane.Debugf("%T %v: Initialize %T %v: %+v", s, s.Name(), step, step, result)

// 	return nil
// }

// func (s *Middleware) eventInitializeWorkflow(params *protocol.InitializeParams, workspace *Workspace, workflow *config.Workflow) error {
// 	for _, step := range workflow.Steps {
// 		switch step.Kind.(type) {
// 		case *config.StepKindLsp:
// 			if err := s.eventInitializeLsp(params, workspace, step); err != nil {
// 				return err
// 			}
// 		default:
// 			err := fmt.Errorf("Initialize unhandled %T %v", step, step)
// 			s.plane.Warnf("%T %v: %v", s, s.Name(), err)
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (s *Middleware) eventInitialize(params *protocol.InitializeParams, workspaces []*Workspace) error {
// 	for _, workspace := range workspaces {
// 		usages := s.plane.Usage().GetFilteredClient(
// 			usage.ConditionFilter{Filetype: nil, Path: &workspace.uri},
// 			usage.EventFilter{Event: util.Ptr(config.EventInitialize)}, s.listenerName, s.clientName)

// 		for _, usage := range usages {
// 			if err := s.eventInitializeWorkflow(params, workspace, usage.Workflow); err != nil {
// 				continue
// 			}
// 			s.plane.Infof("%T %v: Initialize usage %v ok", s, s.Name(), usage.Name)
// 		}
// 	}

// 	return nil
// }

func (s *Middleware) Initialize(params *protocol.InitializeParams, listenerName string, clientName string) (*protocol.InitializeResult, error) {
	s.listenerName = listenerName
	s.clientName = clientName

	for _, workspace := range params.WorkspaceFolders {
		s.workspaceAdd(workspace.URI, workspace.Name)
	}

	// workspaces := s.workspaceFindAll()

	// s.eventInitialize(params, workspaces)

	s.plane.Infof("%T %v: Initialize event done", s, s.Name())

	result := &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    s.textDocumentSyncKind,
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
