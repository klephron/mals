package lsp_tcp

import (
	"context"
	"errors"
	"fmt"
	"mals/internal/listener"
	"mals/internal/plane"
	"net"
)

type ListenerLsp struct {
	listener.Listener
	name  string
	addr  string
	plane plane.Plane
}

func NewListener(name string, port int, plane plane.Plane) (*ListenerLsp, error) {
	l := &ListenerLsp{
		name:  name,
		addr:  fmt.Sprintf(":%d", port),
		plane: plane,
	}
	return l, nil
}

func (s *ListenerLsp) Name() string {
	return s.name
}

func (s *ListenerLsp) Kind() string {
	return Kind()
}

func (s *ListenerLsp) Ipc() string {
	return Ipc()
}

func (s *ListenerLsp) Start(ctx context.Context) error {
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
			s.plane.Log().Errorf("%s: %v", s.Name(), err)
			continue
		}

		s.plane.Log().Infof("%s: accepted %v", s.Name(), conn.RemoteAddr())

		client := newClient(s.plane, conn)

		if err := s.plane.Client().ClientOwn(client, s); err != nil {
			s.plane.Log().Warnf("%s: %v", s.Name(), err)
			continue
		}

		if err := s.plane.Client().ClientServe(client); err != nil {
			s.plane.Log().Warnf("%s: %v", s.Name(), err)
			continue
		}

		s.plane.Log().Infof("%s: started %v", s.Name(), conn.RemoteAddr())
	}
}
