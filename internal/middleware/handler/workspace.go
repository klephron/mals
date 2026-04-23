package handler

import (
	"strings"
)

func (s *Handler) workspaceAdd(uri string, name string) {
	workspace := newWorkspace(uri, name)
	s.workspaces.Store(uri, workspace)

	s.plane.Infof("%T %s: workspace %s added", s, s.Name(), workspace.name)
}

func (s *Handler) workspaceDelete(uri string) {
	workspace, ok := s.workspaces.LoadAndDelete(uri)

	if !ok {
		s.plane.Warnf("%T %s: workspace by uri %s is not present", s, s.Name(), uri)
	}
	s.plane.Infof("%T %s: workspace %s deleted", s, s.Name(), workspace.name)
}

func (s *Handler) workspaceFindAll() []*workspace {
	workspaces := make([]*workspace, 0)

	s.workspaces.Range(func(key string, value *workspace) bool {
		workspaces = append(workspaces, value)
		return true
	})

	return workspaces
}

func (s *Handler) workspaceFindAllByPrefix(uri string) []*workspace {
	workspaces := make([]*workspace, 0)

	s.workspaces.Range(func(key string, value *workspace) bool {
		if strings.HasPrefix(uri, key) {
			workspaces = append(workspaces, value)
		}
		return true
	})

	return workspaces
}
