package lsp

import (
	"context"
	"errors"
	"fmt"
	"mals/internal/listener"
	client "mals/internal/lsp/client/net"
	"mals/internal/plane"
	"net"
)

type ListenerLspTcp struct {
	listener.Listener
	name  string
	addr  string
	plane plane.Plane
}

func New(name string, port int, plane plane.Plane) (*ListenerLspTcp, error) {
	l := &ListenerLspTcp{
		name:  name,
		addr:  fmt.Sprintf(":%d", port),
		plane: plane,
	}
	return l, nil
}

func Kind() string {
	return "lsp"
}

func Ipc() string {
	return "tcp"
}

func (s *ListenerLspTcp) Name() string {
	return s.name
}

func (s *ListenerLspTcp) Kind() string {
	return Kind()
}

func (s *ListenerLspTcp) Ipc() string {
	return Ipc()
}

func (s *ListenerLspTcp) Listen(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)

	if err != nil {
		s.plane.Log().Errorf("%s: %v", s.Name(), err)
		return err
	}

	defer listener.Close()

	s.plane.Log().Infof("%s: listen", s.Name())

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()

		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				s.plane.Log().Infof("%s: closed", s.Name())
				return nil
			}
			s.plane.Log().Warnf("%s: %v", s.Name(), err)
			continue
		}

		s.plane.Log().Infof("%s: accepted %v", s.Name(), conn.RemoteAddr())

		client := client.New(s.plane)
		client.Bind(conn)
	}
}
