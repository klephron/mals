package middleware

import (
	"mals/internal/middleware/document"
	"strings"

	"github.com/puzpuzpuz/xsync/v4"
)

type Workspace struct {
	uri       string
	name      string
	documents *xsync.Map[string, *document.Document]
}

func newWorkspace(uri string, name string) *Workspace {
	return &Workspace{
		uri:       uri,
		name:      name,
		documents: xsync.NewMap[string, *document.Document](),
	}
}

func (s *Middleware) workspaceAdd(uri string, name string) {
	workspace := newWorkspace(uri, name)
	s.workspaces.Store(uri, workspace)

	s.plane.Infof("%s: workspace %s added", s.Name(), workspace.name)
}

func (s *Middleware) workspaceDelete(uri string) {
	workspace, ok := s.workspaces.LoadAndDelete(uri)

	if !ok {
		s.plane.Warnf("%s: workspace by uri %s is not present", s.Name(), uri)
	}
	s.plane.Infof("%s: workspace %s deleted", s.Name(), workspace.name)
}

func (s *Middleware) workspaceFindAll(uri string) []*Workspace {
	workspaces := make([]*Workspace, 0)

	s.workspaces.Range(func(key string, value *Workspace) bool {
		if strings.HasPrefix(uri, key) {
			workspaces = append(workspaces, value)
		}
		return true
	})

	return workspaces
}
