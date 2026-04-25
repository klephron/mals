package handler

import (
	"fmt"
	"mals/internal/middleware/document"
	"mals/internal/plane"
	"mals/internal/scope"
	"mals/pkg/config"

	"github.com/puzpuzpuz/xsync/v4"
)

type workspace struct {
	uri       string
	name      string
	documents *xsync.Map[string, *document.Document]
}

type Handler struct {
	endpoints *config.HandlerLspEndpoints
	resources *xsync.Map[string, *config.HandlerLspResource]

	plane plane.Plane

	listenerName string
	clientName   string
	handlerName  string

	workspaces *xsync.Map[string, *workspace]
}

func newWorkspace(uri string, name string) *workspace {
	return &workspace{
		uri:       uri,
		name:      name,
		documents: xsync.NewMap[string, *document.Document](),
	}
}

func New(listenerName string, clientName string, handlerName string, resources []*config.HandlerLspResource, endpoints *config.HandlerLspEndpoints, plane plane.Plane) *Handler {
	s := Handler{
		endpoints:    endpoints,
		resources:    xsync.NewMap[string, *config.HandlerLspResource](),
		plane:        plane,
		listenerName: listenerName,
		clientName:   clientName,
		handlerName:  handlerName,
		workspaces:   xsync.NewMap[string, *workspace](),
	}

	for _, resource := range resources {
		s.resources.Store(resource.Name, resource)
	}

	return &s
}

func (s *Handler) Name() string {
	return fmt.Sprintf("%v:%v", s.listenerName, s.clientName)
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
