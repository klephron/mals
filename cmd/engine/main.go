package main

import (
	"flag"
	"fmt"
	"log"
	"mals-engine/internal/client"
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

func (p *Params) parse() {
	flag.IntVar(&p.flagPort, "p", 9200, "port to serve")

	flag.Parse()
}

func (e *Engine) setupLogger() {
	e.logger = log.New(os.Stdout, "", log.LUTC|log.Lshortfile|log.Ldate|log.Ltime)
}

func (e *Engine) listen() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", e.flagPort))
	if err != nil {
		return err
	}
	e.logger.Printf("info: listening on port %d", e.flagPort)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			e.logger.Printf("error: %s", err)
			continue
		}

		go func() {
			client := client.NewClient(e.logger, conn)

			defer func() {
				if err := client.Close(); err != nil {
					e.logger.Printf("error: %s", err)
				}
			}()

			client.Listen()
		}()
	}
}

func main() {
	var engine Engine

	engine.parse()
	engine.setupLogger()

	if err := engine.listen(); err != nil {
		engine.logger.Fatal(err)
	}
}
