package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"mals-engine/internal/jsonrpc"
	"mals-engine/internal/logger"
	"net"
)

type Context struct {
	flagPort int

	logger *log.Logger
}

func (ctx Context) connHandle(conn net.Conn) {
	defer func() {
		conn.Close()
		ctx.logger.Printf("info: client disconnected %s", conn.RemoteAddr())
	}()

	ctx.logger.Printf("info: client connected %s", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	scanner.Split(jsonrpc.ScannerSplit)

	for scanner.Scan() {
		bytes := scanner.Bytes()

		msg, _, err := jsonrpc.DecodeRequestMessage(bytes)

		if err != nil {
			ctx.logger.Printf("error: unable to decode %s", err)
			continue
		}

		ctx.logger.Printf("info: handling method %s", msg.Method)

		switch msg.Method {
		case "initialize":
			break
		default:
			ctx.logger.Printf("warn: unhandled method %s", msg.Method)
			break
		}
	}
}

func main() {
	var ctx Context

	flag.IntVar(&ctx.flagPort, "p", 9200, "port to serve")

	flag.Parse()

	logger, err := logger.GetLogger()
	if err != nil {
		panic(err)
	}
	ctx.logger = logger

	ctx.logger.Println("info: started")

	listener, err := net.Listen("tcp", ":"+fmt.Sprint(ctx.flagPort))
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	ctx.logger.Println("info: server listening on port " + fmt.Sprint(ctx.flagPort))

	for {
		conn, err := listener.Accept()
		if err != nil {
			ctx.logger.Printf("error: %s", err)
			continue
		}
		go ctx.connHandle(conn)
	}
}
