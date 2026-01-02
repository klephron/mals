package middleware

import (
	"mals/internal/client"
	"mals/internal/plane"

	"github.com/puzpuzpuz/xsync/v4"
)

type Middleware struct {
	plane       plane.Plane
	client      client.Client
	initialized bool
	workspaces  *xsync.Map[string, *Workspace]
}

func New(plane plane.Plane) *Middleware {
	return &Middleware{
		plane:       plane,
		client:      nil,
		initialized: false,
		workspaces:  xsync.NewMap[string, *Workspace](),
	}
}

func (s *Middleware) Name() string {
	if s.client != nil {
		return s.client.Name()
	}
	return "middleware"
}
