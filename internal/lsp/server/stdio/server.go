package stdio

import (
	"context"
	"mals/internal/lsp/server"
	"mals/internal/plane"
	"mals/pkg/config"
)

type LspServerStdio struct {
	server.LspServer

	name string

	cmd []string

	plane plane.Plane
}

func New(name string, settings *config.LspSettingsStdio, plane plane.Plane) *LspServerStdio {
	cmd := settings.Cmd

	return &LspServerStdio{
		name:  name,
		cmd:   cmd,
		plane: plane,
	}
}

func (s *LspServerStdio) Name() string {
	return s.name
}

func (s *LspServerStdio) Kind() string {
	var settings config.LspSettingsStdio
	return settings.Kind()
}

func (s *LspServerStdio) Run(ctx context.Context) error {
	s.plane.Infof("%T %v: started", s, s.Name())

	<-ctx.Done()

	s.plane.Infof("%T %v: done", s, s.Name())

	return nil
}
