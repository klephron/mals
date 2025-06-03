package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"mals-engine/internal/jsonrpc"
	"net"
	"os"
)

type Params struct {
	flagPort int
}

type Engine struct {
	Params

	logger *log.Logger
}

func (self *Params) parse() {
	flag.IntVar(&self.flagPort, "p", 9200, "port to serve")

	flag.Parse()
}

func (self *Engine) setupLogger() {
	self.logger = log.New(os.Stdout, "", log.LUTC|log.Lshortfile|log.Ldate|log.Ltime)
}

func (self *Engine) listen() error {
	listener, err := net.Listen("tcp", ":"+fmt.Sprint(self.flagPort))
	if err != nil {
		return err
	}
	self.logger.Println("info: listening on port " + fmt.Sprint(self.flagPort))

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			self.logger.Printf("error: %s", err)
			continue
		}

		go func() {
			self.logger.Printf("info: %s connected", conn.RemoteAddr())

			defer func() {
				self.logger.Printf("info: %s disconnected", conn.RemoteAddr())
				conn.Close()
			}()

			self.connHandle(conn)
		}()
	}
}

func (self *Engine) connHandle(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	scanner.Split(jsonrpc.ScannerSplit)

	for scanner.Scan() {
		bytes := scanner.Bytes()

		msg, _, err := jsonrpc.DecodeRequestMessage(bytes)

		if err != nil {
			self.logger.Printf("error: %s unable to decode %s", conn.RemoteAddr(), err)
			continue
		}

		self.logger.Printf("info: %s handling method %s", conn.RemoteAddr(), msg.Method)

		switch msg.Method {
		case "initialize":
			break
		default:
			self.logger.Printf("warn: %s unhandled method %s", conn.RemoteAddr(), msg.Method)
			break
		}
	}
}

func main() {
	var engine Engine

	engine.parse()
	engine.setupLogger()

	engine.listen()
}
