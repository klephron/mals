package middleware

import (
	"fmt"
	"mals/internal/scope"
	"mals/internal/usage"
	"mals/internal/util"
	"mals/pkg/config"
)

func (s *Middleware) eventShutdownLsp(_ *Workspace, step *config.Step) error {
	if step.Scope != "client" {
		s.plane.Warnf("Shutdown %T %v scope %v unsupported, set to client", step, step, step.Scope)
	}
	scope := scope.NewScopeClient(s.listenerName, s.clientName)

	lspName := step.Kind.(*config.StepKindLsp).Name

	lspKey, token, err := s.plane.Scope().LspAcquire(lspName, scope)
	if err != nil {
		s.plane.Errorf("Shutdown %T %v: %v", step, step, err)
		return err
	}
	defer s.plane.Scope().LspRelease(lspKey, token)

	err = s.plane.Lsp().EventShutdown(lspKey)
	if err != nil {
		s.plane.Errorf("Shutdown %T %v: %v", step, step, err)
		return nil
	}

	s.plane.Debugf("Shutdown %T %v", step, step)

	return nil
}

func (s *Middleware) eventShutdownWorkflow(workspace *Workspace, workflow *config.Workflow) error {
	for _, step := range workflow.Steps {
		switch step.Kind.(type) {
		case *config.StepKindLsp:
			if err := s.eventShutdownLsp(workspace, step); err != nil {
				return err
			}
		default:
			err := fmt.Errorf("Shutdown unhandled %T %v", step, step)
			s.plane.Warnf("%v", err)
			return err
		}
	}

	return nil
}

func (s *Middleware) eventShutdown(workspaces []*Workspace) error {
	for _, workspace := range workspaces {
		usages := s.plane.Usage().GetFilteredClient(
			usage.ConditionFilter{Filetype: nil, Path: &workspace.uri},
			usage.EventFilter{Event: util.Ptr(config.EventShutdown)}, s.listenerName, s.clientName)

		for _, usage := range usages {
			if err := s.eventShutdownWorkflow(workspace, usage.Workflow); err != nil {
				continue
			}
			s.plane.Infof("Shutdown %T %v: usage %v ok", s, s.Name(), usage.Name)
		}
	}

	return nil
}

func (s *Middleware) Shutdown() error {
	workspaces := s.workspaceFindAll()

	s.eventShutdown(workspaces)

	s.workspaces.Range(func(key string, value *Workspace) bool {
		s.workspaceDelete(key)
		return true
	})

	s.initialized = false

	return nil
}
