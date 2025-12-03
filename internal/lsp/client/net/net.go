package client

import (
	"bufio"
	"fmt"
	"mals/internal/control/controller"
	"mals/internal/lsp/client/generic"
	"net"
)

type ClientNet struct {
	client.ClientGeneric
	conn net.Conn
}

func New(controller *controller.Controller) *ClientNet {
	return &ClientNet{
		ClientGeneric: *client.New(controller),
		conn:          nil,
	}
}

func (s *ClientNet) Name() string {
	return fmt.Sprintf("client[%s]", s.conn.RemoteAddr())
}

func (s *ClientNet) Bind(conn net.Conn) {
	s.ClientGeneric.Bind(bufio.NewScanner(conn), bufio.NewWriter(conn))
	s.conn = conn
}

func (s *ClientNet) Unbind() error {
	if err := s.ClientGeneric.Unbind(); err != nil {
		return err
	}
	return s.conn.Close()
}

func (s *ClientNet) Close() error {
	if err := s.Unbind(); err != nil {
		return err
	}
	return s.ClientGeneric.Close()
}
