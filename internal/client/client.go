package client

import (
	"bufio"
	"context"
	"log"
	"mals-engine/internal/jsonrpc"
	"mals-engine/internal/workspace"
	"net"
)

type Client struct {
	logger     *log.Logger
	conn       net.Conn
	scanner    *bufio.Scanner
	writer     *bufio.Writer
	workspaces map[string]*workspace.Workspace // path should be cleaned
}

func NewClient(logger *log.Logger, conn net.Conn) (c *Client) {
	c = &Client{
		logger:     logger,
		conn:       conn,
		scanner:    bufio.NewScanner(conn),
		writer:     bufio.NewWriter(conn),
		workspaces: make(map[string]*workspace.Workspace),
	}
	c.scanner.Split(jsonrpc.ScannerSplit)
	return
}

func (c *Client) Serve(ctx context.Context) {
	defer c.Close()

	c.LogInfoPrintf("listening")

	bytesC := make(chan []byte)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(bytesC)
				return
			default:
				if !c.scanner.Scan() {
					close(bytesC)
					return
				}
				bytesC <- c.scanner.Bytes()
			}
		}
	}()

	for {
		select {
		case <-ctx.Done(): // can wait on scan
			return
		case bytes, ok := <-bytesC:
			if !ok {
				return
			}
			c.HandleLspRequest(bytes)
		}
	}
}

func (c *Client) Close() error {
	c.LogInfoPrintf("close")
	if err := c.conn.Close(); err != nil {
		return err
	}
	return nil
}
