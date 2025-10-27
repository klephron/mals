package client

import (
	"bufio"
	"context"
	"log"
	"mals/internal/jsonrpc"
	"mals/internal/model"
	"mals/internal/workspace"
	"net"
)

type WorkspaceConfig struct {
	DefaultModel model.ModelService
}

type Config struct {
	Workspace WorkspaceConfig
}

type Client struct {
	logger     *log.Logger
	conn       net.Conn
	scanner    *bufio.Scanner
	writer     *bufio.Writer
	workspaces map[string]*workspace.Workspace // path should be cleaned
	config     Config
}

func NewClient(logger *log.Logger, conn net.Conn, config Config) (c *Client) {
	c = &Client{
		logger:     logger,
		conn:       conn,
		scanner:    bufio.NewScanner(conn),
		writer:     bufio.NewWriter(conn),
		workspaces: make(map[string]*workspace.Workspace),
		config:     config,
	}
	c.scanner.Split(jsonrpc.ScannerSplit)
	return
}

func (c *Client) Serve(ctx context.Context) {
	defer c.Close()

	bytesC := make(chan []byte)
	defer close(bytesC)

	c.LogInfoPrintf("listening")

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if !c.scanner.Scan() {
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
