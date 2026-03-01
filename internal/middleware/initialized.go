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
		s.plane.Warnf("%T %v: Initialized %T %v scope %v unsupported, set to client", s, s.Name(), step, step, step.Scope)
	}
	scope := scope.NewScopeClient(s.listenerName, s.clientName)

	lspName := step.Kind.(*config.StepKindLsp).Name

	lspKey, token, err := s.plane.Scope().LspAcquire(lspName, scope)
	if err != nil {
		s.plane.Errorf("%T %v: Initialized %T %v: %v", s, s.Name(), step, step, err)
		return err
	}
	defer s.plane.Scope().LspRelease(lspKey, token)

	lspParams := &protocol.InitializedParams{}

	err = s.plane.Lsp().EventInitialized(lspKey, lspParams)
	if err != nil {
		s.plane.Errorf("%T %v: Initialized %T %v: %v", s, s.Name(), step, step, err)
		return nil
	}

	s.plane.Debugf("%T %v: Initialized %T %v", s, s.Name(), step, step)

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
			s.plane.Warnf("%T %v: %v", s, s.Name(), err)
			return err
		}
	}

	return nil
}

func (s *Middleware) eventInitialized(params *protocol.InitializedParams, workspaces []*Workspace) error {
	for _, workspace := range workspaces {
		usages := s.plane.Usage().GetFilteredClient(
			usage.ConditionFilter{Filetype: nil, Path: &workspace.uri},
			usage.EventFilter{Event: util.Ptr(config.EventInitialized)}, s.listenerName, s.clientName)

		for _, usage := range usages {
			if err := s.eventInitializedWorkflow(params, workspace, usage.Workflow); err != nil {
				continue
			}
			s.plane.Infof("%T %v: Initialized usage %v ok", s, s.Name(), usage.Name)
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

	s.plane.Infof("%T %v: Initialized event done", s, s.Name())

	return nil
}
