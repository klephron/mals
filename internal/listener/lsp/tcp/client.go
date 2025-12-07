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
	client.Client = client
	return client
}

func (s *ClientLspTcp) Name() string {
	if s.conn != nil {
		return fmt.Sprintf("%s", s.conn.RemoteAddr())
	}
	return fmt.Sprintf("%v", nil)
}

func (s *ClientLspTcp) Close() error {
	if err := s.ClientLsp.Close(); err != nil {
		return err
	}
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			return err
		}
		s.conn = nil
	}
	return nil
}

func (s *ClientLspTcp) Bind(conn net.Conn) error {
	s.Unbind()
	s.conn = conn
	return s.ClientLsp.Bind(bufio.NewScanner(conn), bufio.NewWriter(conn))
}

func (s *ClientLspTcp) Unbind() error {
	if err := s.ClientLsp.Unbind(); err != nil {
		return err
	}
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			return err
		}
		s.conn = nil
	}
	return nil
}
