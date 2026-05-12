package handler

import (
	"mals/pkg/config"
	"mals/pkg/info"
	"mals/third_party/lsp"
	"os"
)

func (s *Handler) InitializeDefault(params *lsp.InitializeParams) error {
	for _, workspace := range params.WorkspaceFolders {
		s.workspaceAdd(workspace.URI, workspace.Name)
	}

	s.resources.Range(func(key string, value *config.HandlerLspResource) bool {
		scope, err := s.getResourceScope(value.Scope)
		if err != nil {
			s.plane.Errorf("%T %v: Initialize %v", s, s.Name(), err)
			return true
		}

		switch vs := value.Spec.(type) {
		case *config.HandlerLspResourceSpecLsp:
			lspName, token, err := s.plane.Scope().LspAcquire(vs.Name, scope)
			if err != nil {
				s.plane.Errorf("%T %v: Initialize %v", s, s.Name(), err)
				return true
			}

			defer func() {
				if err := s.plane.Scope().LspRelease(lspName, token); err != nil {
					s.plane.Errorf("%T %v: Initialize %T %v: %v", s, s.Name(), err)
				}
			}()

			workspaces := s.workspaceFindAll()
			workspaceFolders := make([]lsp.WorkspaceFolder, len(workspaces))
			for i, workspace := range workspaces {
				workspaceFolders[i] = lsp.WorkspaceFolder{
					URI:  workspace.Uri(),
					Name: workspace.Name(),
				}
			}

			lspParams := &lsp.InitializeParams{
				XInitializeParams: lsp.XInitializeParams{
					ProcessID: int32(os.Getpid()),
					ClientInfo: &lsp.ClientInfo{
						Name:    info.MiddlewareClientName,
						Version: info.MiddlewareVersion,
					},
					Locale: params.Locale,
					Capabilities: lsp.ClientCapabilities{
						TextDocument: lsp.TextDocumentClientCapabilities{
							Synchronization: &lsp.TextDocumentSyncClientCapabilities{
								DynamicRegistration: false,
								WillSave:            false,
								WillSaveWaitUntil:   false,
								DidSave:             false,
							},
							Completion: lsp.CompletionClientCapabilities{
								DynamicRegistration: false,
							},
						},
					},
					InitializationOptions: params.InitializationOptions,
					Trace:                 params.Trace,
				},
				WorkspaceFoldersInitializeParams: lsp.WorkspaceFoldersInitializeParams{
					WorkspaceFolders: workspaceFolders,
				},
			}

			if _, err := s.plane.Lsp().Initialize(lspName, lspParams); err != nil {
				s.plane.Errorf("%T %v: Initialize %T %v: %v", s, s.Name(), err)
			}

		case *config.HandlerLspResourceSpecModel:
			// trigger acquisition of resource
			modelName, token, err := s.plane.Scope().ModelAcquire(vs.Name, scope)
			if err != nil {
				s.plane.Errorf("%T %v: Initialize %T %v: %v", s, s.Name(), err)
				return true
			}
			if err := s.plane.Scope().ModelRelease(modelName, token); err != nil {
				s.plane.Errorf("%T %v: Initialize %v", s, s.Name(), err)
				return true
			}
		default:
			s.plane.Errorf("%T %v: Initialize unexpected spec %T", s, s.Name(), vs)
		}

		return true
	})

	return nil
}

func (s *Handler) Initialize(params *lsp.InitializeParams) error {
	if *s.endpoints.Initialize.Default {
		err := s.InitializeDefault(params)
		if err != nil {
			return err
		}
	}

	s.plane.Infof("%T %v: Initialize done", s, s.Name())

	return nil
}
