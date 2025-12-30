package lsp_tcp

import (
	"bufio"
	"context"
	"fmt"
	"mals/internal/lsp/client"
	"mals/internal/plane"
	"net"
)

type ClientLspTcp struct {
	client.ClientLsp
	conn net.Conn
}

func newClient(plane plane.Plane, conn net.Conn) *ClientLspTcp {
	c := &ClientLspTcp{
		ClientLsp: *client.New(plane, bufio.NewScanner(conn), bufio.NewWriter(conn)),
		conn:      conn,
	}
	c.Client = c
	return c
}

func (s *ClientLspTcp) Name() string {
	return fmt.Sprintf("%s", s.conn.RemoteAddr())
}

func (s *ClientLspTcp) Serve(ctx context.Context) (err error) {
	err = s.ClientLsp.Serve(ctx)

	if err != nil {
		s.conn.Close()
	} else {
		err = s.conn.Close()
	}

	return
}
