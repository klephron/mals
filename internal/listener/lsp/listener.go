package lsp

import (
	"context"
	"errors"
	"fmt"
	"mals/internal/plane"
	"net"
)

type ListenerLsp struct {
	name  string
	addr  string
	plane plane.Plane
}

func NewTcp(name string, port int, plane plane.Plane) (*ListenerLsp, error) {
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

func (s *ListenerLsp) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)

	if err != nil {
		s.plane.Errorf("%s: %v", s.Name(), err)
		return err
	}

	defer listener.Close()

	s.plane.Infof("%s: listen", s.Name())

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()

		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				s.plane.Infof("%s: closed", s.Name())
				return nil
			}
			s.plane.Errorf("%s: %v", s.Name(), err)
			continue
		}

		s.plane.Infof("%s: accepted %v", s.Name(), conn.RemoteAddr())

		client := newLspClient(s.plane, s.name, conn)

		if err := s.plane.Listener().LspClientOwn(s.Name(), client); err != nil {
			s.plane.Warnf("%s: %v", s.Name(), err)
			continue
		}

		if err := s.plane.Listener().LspClientServe(s.Name(), client.Name()); err != nil {
			s.plane.Warnf("%s: %v", s.Name(), err)
			continue
		}

		s.plane.Infof("%s: started %v", s.Name(), conn.RemoteAddr())
	}
}
