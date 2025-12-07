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

func NewClient(plane plane.Plane, conn net.Conn) *ClientLspTcp {
	client := &ClientLspTcp{
		ClientLsp: *lsp.New(plane, bufio.NewScanner(conn), bufio.NewWriter(conn)),
		conn:      conn,
	}
	client.Client = client
	return client
}

func (s *ClientLspTcp) Name() string {
	return fmt.Sprintf("%s", s.conn.RemoteAddr())
}

func (s *ClientLspTcp) Close() error {
	if err := s.ClientLsp.Close(); err != nil {
		return err
	}
	if err := s.conn.Close(); err != nil {
		return err
	}
	return nil
}
