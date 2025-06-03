package client

import (
	"bufio"
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

func (c *Client) Listen() {
	c.logger.Printf("info: %s listening", c.conn.RemoteAddr())

	scanner := bufio.NewScanner(c.conn)
	scanner.Split(jsonrpc.ScannerSplit)

	for scanner.Scan() {
		bytes := scanner.Bytes()

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
	if err := scanner.Err(); err != nil {
		c.logger.Printf("error: %s scanner error: %v", c.conn.RemoteAddr(), err)
	}
}

func (c *Client) Close() error {
	c.logger.Printf("info: %s disconnected", c.conn.RemoteAddr())
	return c.conn.Close()
}
