package middleware

import (
	"fmt"
	"mals/internal/lsp/protocol"
	"mals/internal/scope"
	"mals/internal/usage"
	"mals/internal/util"
	"mals/pkg/config"
)

func (s *Middleware) eventInitializedLsp(_ *protocol.InitializedParams, _ *Workspace, step *config.Step) error {
	if step.Scope != "client" {
		s.plane.Warnf("Initialized %T %v scope %v unsupported, set to client", step, step, step.Scope)
	}
	scope := scope.NewScopeClient(s.client.Name())

	lspName := step.Kind.(*config.StepKindLsp).Name

	lspKey, token, err := s.plane.ScopeLspAcquire(lspName, scope)
	if err != nil {
		s.plane.Errorf("Initialized %T %v: %v", step, step, err)
		return err
	}
	defer s.plane.ScopeLspRelease(lspKey, token)

	lspParams := &protocol.InitializedParams{}

	err = s.plane.LspEventInitialized(lspKey, lspParams)
	if err != nil {
		s.plane.Errorf("Initialized %T %v: %v", step, step, err)
		return nil
	}

	s.plane.Infof("Initialized %T %v", step, step)

	return nil
}

func (s *Middleware) eventInitializedWorkflow(params *protocol.InitializedParams, workspace *Workspace, workflow *config.Workflow) error {
	for _, step := range workflow.Steps {
		switch step.Kind.(type) {
		case *config.StepKindLsp:
			if err := s.eventInitializedLsp(params, workspace, step); err != nil {
				return err
			}
		default:
			err := fmt.Errorf("Initialized unhandled %T %v", step, step)
			s.plane.Warnf("%v", err)
			return err
		}
	}

	return nil
}

func (s *Middleware) eventInitialized(params *protocol.InitializedParams, workspaces []*Workspace) error {
	for _, workspace := range workspaces {
		usages := s.plane.UsageGetFilteredClient(
			usage.ConditionFilter{Filetype: nil, Path: &workspace.uri},
			usage.EventFilter{Event: util.Ptr(config.EventInitialized)}, s.client.Name())

		for _, usage := range usages {
			if err := s.eventInitializedWorkflow(params, workspace, usage.Workflow); err != nil {
				continue
			}
			s.plane.Infof("Initialized %T %v: usage %v ok", s, s.Name(), usage.Name)
		}
	}
	return nil
}

func (s *Middleware) Initialized(params *protocol.InitializedParams) error {
	if s.initialized {
		return fmt.Errorf("%v: already initialized", s.Name())
	}

	s.initialized = true

	workspaces := s.workspaceFindAll()

	s.eventInitialized(params, workspaces)

	s.plane.Infof("Initialized %v: event done", s.Name())

	return nil
}
