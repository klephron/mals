package handler

import (
	"fmt"
	"mals/internal/middleware/document"
	"mals/internal/plane"
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

	workspaces *xsync.Map[string, *workspace]
}

func newWorkspace(uri string, name string) *workspace {
	return &workspace{
		uri:       uri,
		name:      name,
		documents: xsync.NewMap[string, *document.Document](),
	}
}

func New(listenerName string, clientName string, resources []*config.HandlerLspResource, endpoints *config.HandlerLspEndpoints, plane plane.Plane) *Handler {
	s := Handler{
		resources: xsync.NewMap[string, *config.HandlerLspResource](),
		endpoints: endpoints,
		plane:     plane,
	}

	for _, resource := range resources {
		s.resources.Store(resource.Name, resource)
	}

	return &s
}

func (s *Handler) Name() string {
	return fmt.Sprintf("%v:%v", s.listenerName, s.clientName)
}
