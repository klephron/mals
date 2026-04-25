package handler

import (
	"fmt"
	"mals/internal/middleware/workspace"
	"mals/internal/plane"
	"mals/internal/scope"
	"mals/pkg/config"
	"strings"

	"github.com/puzpuzpuz/xsync/v4"
)

type Handler struct {
	endpoints *config.HandlerLspEndpoints
	resources *xsync.Map[string, *config.HandlerLspResource]

	plane plane.Plane

	listenerName string
	clientName   string
	handlerName  string

	workspaces *xsync.Map[string, *workspace.Workspace]
}

func New(listenerName string, clientName string, handlerName string, resources []*config.HandlerLspResource, endpoints *config.HandlerLspEndpoints, plane plane.Plane) *Handler {
	s := Handler{
		endpoints:    endpoints,
		resources:    xsync.NewMap[string, *config.HandlerLspResource](),
		plane:        plane,
		listenerName: listenerName,
		clientName:   clientName,
		handlerName:  handlerName,
		workspaces:   xsync.NewMap[string, *workspace.Workspace](),
	}

	for _, resource := range resources {
		s.resources.Store(resource.Name, resource)
	}

	return &s
}

func (s *Handler) Name() string {
	return fmt.Sprintf("%v:%v", s.listenerName, s.clientName)
}

func (s *Handler) getScope() *scope.Scope {
	return scope.NewScopeHandler(s.listenerName, s.clientName, s.handlerName)
}

func (s *Handler) getResourceScope(resourceScope config.HandlerLspResourceScope) (*scope.Scope, error) {
	switch resourceScope {
	case config.HandlerLspResourceScopeGlobal:
		return scope.NewScopeGlobal(), nil
	case config.HandlerLspResourceScopeClient:
		return scope.NewScopeClient(s.listenerName, s.clientName), nil
	case config.HandlerLspResourceScopeHandler:
		return scope.NewScopeHandler(s.listenerName, s.clientName, s.handlerName), nil
	default:
		return nil, fmt.Errorf("unknown scope kind")
	}
}

func (s *Handler) workspaceAdd(uri string, name string) {
	workspace := workspace.New(uri, name)
	s.workspaces.Store(uri, workspace)

	s.plane.Infof("%T %s: workspace %s added", s, s.Name(), workspace.Name())
}

func (s *Handler) workspaceDelete(uri string) {
	workspace, ok := s.workspaces.LoadAndDelete(uri)

	if !ok {
		s.plane.Warnf("%T %s: workspace by uri %s is not present", s, s.Name(), uri)
	}
	s.plane.Infof("%T %s: workspace %s deleted", s, s.Name(), workspace.Name())
}

func (s *Handler) workspaceFindAll() []*workspace.Workspace {
	workspaces := make([]*workspace.Workspace, 0)

	s.workspaces.Range(func(key string, value *workspace.Workspace) bool {
		workspaces = append(workspaces, value)
		return true
	})

	return workspaces
}

func (s *Handler) workspaceFindAllByPrefix(uri string) []*workspace.Workspace {
	workspaces := make([]*workspace.Workspace, 0)

	s.workspaces.Range(func(key string, value *workspace.Workspace) bool {
		if strings.HasPrefix(uri, key) {
			workspaces = append(workspaces, value)
		}
		return true
	})

	return workspaces
}
