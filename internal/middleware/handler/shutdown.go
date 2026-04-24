package handler

import (
// "fmt"
// "mals/internal/scope"
// "mals/internal/util"
// "mals/pkg/config"
)

// func (s *Middleware) eventShutdownLsp(_ *Workspace, step *config.Step) error {
// 	if step.Scope != "client" {
// 		s.plane.Warnf("%T %v: Shutdown %T %v scope %v unsupported, set to client",
// 			s, s.Name(), step, step, step.Scope)
// 	}
// 	scope := scope.NewScopeClient(s.listenerName, s.clientName)

// 	lspName := step.Kind.(*config.StepKindLsp).Name

// 	lspKey, token, err := s.plane.Scope().LspAcquire(lspName, scope)
// 	if err != nil {
// 		s.plane.Errorf("%T %v: Shutdown %T %v: %v", s, s.Name(), step, step, err)
// 		return err
// 	}
// 	defer s.plane.Scope().LspRelease(lspKey, token)

// 	err = s.plane.Lsp().EventShutdown(lspKey)
// 	if err != nil {
// 		s.plane.Errorf("%T %v: Shutdown %T %v: %v", s, s.Name(), step, step, err)
// 		return nil
// 	}

// 	s.plane.Debugf("%T %v: Shutdown %T %v", s, s.Name(), step, step)

// 	return nil
// }

// func (s *Middleware) eventShutdownWorkflow(workspace *Workspace, workflow *config.Workflow) error {
// 	for _, step := range workflow.Steps {
// 		switch step.Kind.(type) {
// 		case *config.StepKindLsp:
// 			if err := s.eventShutdownLsp(workspace, step); err != nil {
// 				return err
// 			}
// 		default:
// 			err := fmt.Errorf("Shutdown unhandled %T %v", step, step)
// 			s.plane.Warnf("%T %v: %v", s, s.Name(), err)
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (s *Middleware) eventShutdown(workspaces []*Workspace) error {
// 	for _, workspace := range workspaces {
// 		usages := s.plane.Usage().GetFilteredClient(
// 			usage.ConditionFilter{Filetype: nil, Path: &workspace.uri},
// 			usage.EventFilter{Event: util.Ptr(config.EventShutdown)}, s.listenerName, s.clientName)

// 		for _, usage := range usages {
// 			if err := s.eventShutdownWorkflow(workspace, usage.Workflow); err != nil {
// 				continue
// 			}
// 			s.plane.Infof("%T %v: Shutdown usage %v ok", s, s.Name(), usage.Name)
// 		}
// 	}

// 	return nil
// }

func (s *Handler) Shutdown() error {
	// workspaces := s.workspaceFindAll()

	// s.eventShutdown(workspaces)

	// s.workspaces.Range(func(key string, value *workspace) bool {
	// 	s.workspaceDelete(key)
	// 	return true
	// })

	// s.initialized = false

	return nil
}
