package main

import (
	"bufio"
	"flag"
	"log"
	"mals-engine/internal/jsonrpc"
	"mals-engine/internal/logger"
	"mals-engine/pkg/message"
	"os"
)

type Context struct {
	logPath string
	isStdio bool
	logger  *log.Logger
}

func (ctx Context) handle(msg message.RequestMessage) {
	switch msg.Method {
	case "initialize":
		break
	default:
		ctx.logger.Printf("warn: unhandled method %s", msg.Method)
		break
	}
}

func prepare() (ctx Context) {
	flag.StringVar(&ctx.logPath, "l", "lsp.log", "file path to log info")
	flag.BoolVar(&ctx.isStdio, "stdio", true, "use stdio for connection")
	flag.Parse()
	logger, err := logger.GetLogger(ctx.logPath)
	if err != nil {
		panic(err)
	}
	ctx.logger = logger
	return
}

func main() {
	ctx := prepare()

	ctx.logger.Println("info: started")

	if !ctx.isStdio {
		panic("adapter supports only stdio")
	}

	if ctx.isStdio {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(jsonrpc.ScannerSplit)

		// writer := os.Stdout

		for scanner.Scan() {
			msg := scanner.Bytes()

			message, _, err := jsonrpc.DecodeRequestMessage(msg)

			if err != nil {
				ctx.logger.Printf("error: unable to decode %s", err)
				continue
			}

			ctx.handle(message)
		}
	}
}
