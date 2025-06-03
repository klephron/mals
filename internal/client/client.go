package client

import (
	"bufio"
	"context"
	"log"
	"mals-engine/internal/jsonrpc"
	"net"
)

type Client struct {
	logger *log.Logger

	conn net.Conn
}

func NewClient(logger *log.Logger, conn net.Conn) *Client {
	return &Client{logger: logger, conn: conn}
}

func (c *Client) Serve(ctx context.Context) {
	defer c.Close()

	c.logger.Printf("info: %s listening", c.conn.RemoteAddr())

	scanner := bufio.NewScanner(c.conn)
	scanner.Split(jsonrpc.ScannerSplit)

	bytesC := make(chan []byte)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(bytesC)
				return
			default:
				if !scanner.Scan() {
					close(bytesC)
					return
				}
				bytesC <- scanner.Bytes()
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
			msg, _, err := jsonrpc.DecodeRequestMessage(bytes)

			if err != nil {
				c.logger.Printf("error: %s unable to decode %s", c.conn.RemoteAddr(), err)
				continue
			}

			c.logger.Printf("info: %s handling method %s", c.conn.RemoteAddr(), msg.Method)

			switch msg.Method {
			case "initialize":
				break
			default:
				c.logger.Printf("warn: %s unhandled method %s", c.conn.RemoteAddr(), msg.Method)
				break
			}
		}
	}
}

func (c *Client) Close() error {
	c.logger.Printf("info: %s close", c.conn.RemoteAddr())
	if err := c.conn.Close(); err != nil {
		return err
	}
	return nil
}
