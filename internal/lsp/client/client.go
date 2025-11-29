package client

import (
	"bufio"
	"context"
	"fmt"
	"mals/internal/jsonrpc"
	"mals/internal/state"
	"net"
)

type Client struct {
	state   state.State
	conn    net.Conn
	scanner *bufio.Scanner
	writer  *bufio.Writer
}

func New(state state.State, conn net.Conn) (c *Client) {
	c = &Client{
		state:   state,
		conn:    conn,
		scanner: bufio.NewScanner(conn),
		writer:  bufio.NewWriter(conn),
	}
	c.scanner.Split(jsonrpc.ScannerSplit)
	return
}

func (s *Client) Scan(ctx context.Context, ch chan<- []byte) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !s.scanner.Scan() {
				return
			}
			ch <- s.scanner.Bytes()
			s.state.Debug(fmt.Sprintf("%s: scanned %s", s.logPrefix(), string(s.scanner.Bytes())))
		}
	}
}

func (s *Client) Listen(ctx context.Context) {
	ch := make(chan []byte)
	defer close(ch)

	s.state.Info(fmt.Sprintf("%s: listening", s.logPrefix()))

	scanCtx, scanCancel := context.WithCancel(ctx)

	go func() {
		s.Scan(scanCtx, ch)
		scanCancel()
	}()

	for {
		select {
		case <-scanCtx.Done():
			s.state.Info(fmt.Sprintf("%s: done", s.logPrefix()))
			return
		case bytes := <-ch:
			s.LspHandle(bytes)
		}
	}
}

func (s *Client) Close() error {
	if err := s.conn.Close(); err != nil {
		s.state.Error(fmt.Sprintf("%s: %v", s.logPrefix(), err))
		return err
	}
	s.state.Info(fmt.Sprintf("%s: closed", s.logPrefix()))
	return nil
}

func (s *Client) logPrefix() string {
	return fmt.Sprintf("client[%s]", s.conn.RemoteAddr())
}
