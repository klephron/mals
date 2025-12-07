package apitcp

import (
	"context"
	"fmt"
	"mals/internal/listener"
	"mals/internal/plane"
)

type ListenerApiTcp struct {
	listener.Listener
	name  string
	addr  string
	plane plane.Plane
}

func NewListener(name string, port int, plane plane.Plane) (*ListenerApiTcp, error) {
	l := &ListenerApiTcp{
		name:  name,
		addr:  fmt.Sprintf(":%d", port),
		plane: plane,
	}
	return l, nil
}

func (s *ListenerApiTcp) Name() string {
	return s.name
}

func (s *ListenerApiTcp) Kind() string {
	return Kind()
}

func (s *ListenerApiTcp) Ipc() string {
	return Ipc()
}

func (s *ListenerApiTcp) Listen(ctx context.Context) error {
	return nil
}
