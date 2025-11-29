package client

import (
	"bufio"
	"context"
	"fmt"
	"mals/internal/jsonrpc"
	// "mals/internal/lsp/workspace"
	"mals/internal/state"
	"net"
)

type Client struct {
	state   state.State
	conn    net.Conn
	scanner *bufio.Scanner
	writer  *bufio.Writer
	// workspace *workspace.Workspace
}

func New(state state.State, conn net.Conn) (c *Client) {
	c = &Client{
		state:   state,
		conn:    conn,
		scanner: bufio.NewScanner(conn),
		writer:  bufio.NewWriter(conn),
		// workspace: nil,
	}
	c.scanner.Split(jsonrpc.ScannerSplit)
	return
}

func (s *Client) Serve(ctx context.Context) {
	defer s.Close()

	bytesC := make(chan []byte)
	defer close(bytesC)

	s.state.Info(fmt.Sprintf("%s: serving", s.logPrefix()))

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if !s.scanner.Scan() {
					return
				}
				bytesC <- s.scanner.Bytes()
				s.state.Debug(fmt.Sprintf("%s: scanned %s", string(s.scanner.Bytes())))
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-bytesC:
			if !ok {
				return
			}
			// s.LspHandle(bytes)
		}
	}
}

func (s *Client) Close() error {
	if err := s.conn.Close(); err != nil {
		s.state.Error(fmt.Sprintf("%s: %v", err))
		return err
	}
	s.state.Info(fmt.Sprintf("%s: closed", s.logPrefix()))
	return nil
}

func (s *Client) logPrefix() string {
	return fmt.Sprintf("client[%s]", s.conn.RemoteAddr())
}
