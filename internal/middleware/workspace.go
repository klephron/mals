package middleware

import (
	"mals/internal/middleware/document"
	"strings"

	"github.com/puzpuzpuz/xsync/v4"
)

func newWorkspace(uri string, name string) *workspace {
	return &workspace{
		uri:       uri,
		name:      name,
		documents: xsync.NewMap[string, *document.Document](),
	}
}

func (s *Middleware) workspaceAdd(uri string, name string) {
	workspace := newWorkspace(uri, name)
	s.workspaces.Store(uri, workspace)

	s.plane.Infof("%T %s: workspace %s added", s, s.Name(), workspace.name)
}

func (s *Middleware) workspaceDelete(uri string) {
	workspace, ok := s.workspaces.LoadAndDelete(uri)

	if !ok {
		s.plane.Warnf("%T %s: workspace by uri %s is not present", s, s.Name(), uri)
	}
	s.plane.Infof("%T %s: workspace %s deleted", s, s.Name(), workspace.name)
}

func (s *Middleware) workspaceFindAll() []*workspace {
	workspaces := make([]*workspace, 0)

	s.workspaces.Range(func(key string, value *workspace) bool {
		workspaces = append(workspaces, value)
		return true
	})

	return workspaces
}

func (s *Middleware) workspaceFindAllByPrefix(uri string) []*workspace {
	workspaces := make([]*workspace, 0)

	s.workspaces.Range(func(key string, value *workspace) bool {
		if strings.HasPrefix(uri, key) {
			workspaces = append(workspaces, value)
		}
		return true
	})

	return workspaces
}
