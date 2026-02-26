package lsp

import (
	"bufio"
	"context"
	"fmt"
	"mals/internal/lsp/client"
	"mals/internal/plane"
	"net"
)

type ListenerLspClient struct {
	conn   net.Conn
	client *client.LspClient
}

func newLspClient(plane plane.Plane, listenerName string, conn net.Conn) *ListenerLspClient {
	s := &ListenerLspClient{
		conn:   conn,
		client: nil,
	}

	s.client = client.New(listenerName, fmt.Sprintf("%s", s.conn.RemoteAddr()),
		plane, bufio.NewScanner(conn), bufio.NewWriter(conn))

	return s
}

func (s *ListenerLspClient) Name() string {
	return s.client.Name()
}

func (s *ListenerLspClient) Serve(ctx context.Context) (err error) {
	err = s.client.Serve(ctx)

	if err != nil {
		s.conn.Close()
	} else {
		err = s.conn.Close()
	}

	return
}
