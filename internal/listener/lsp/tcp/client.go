package lsptcp

import (
	"bufio"
	"fmt"
	"mals/internal/client/lsp"
	"mals/internal/plane"
	"net"
)

type ClientLspTcp struct {
	lsp.ClientLsp
	conn net.Conn
}

func NewClient(plane plane.Plane) *ClientLspTcp {
	client := &ClientLspTcp{
		ClientLsp: *lsp.New(plane),
		conn:      nil,
	}
	client.ClientLsp.Client = client
	return client
}

func (s *ClientLspTcp) Name() string {
	return fmt.Sprintf("%s", s.conn.RemoteAddr())
}

func (s *ClientLspTcp) Bind(conn net.Conn) error {
	if s.conn != nil {
		if err := s.Unbind(); err != nil {
			return err
		}
	}
	s.conn = conn
	return s.ClientLsp.Bind(bufio.NewScanner(conn), bufio.NewWriter(conn))
}

func (s *ClientLspTcp) Unbind() error {
	if err := s.ClientLsp.Unbind(); err != nil {
		return err
	}
	return s.conn.Close()
}

func (s *ClientLspTcp) Close() error {
	if err := s.Unbind(); err != nil {
		return err
	}
	return nil
}
