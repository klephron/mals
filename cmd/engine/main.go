package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mals-engine/internal/client"
	"mals-engine/pkg/config"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Params struct {
	flagPort       int
	flagConfigPath string
}

type Engine struct {
	Params
	logger  *log.Logger
	config  *config.Config
	clients sync.Map
}

func (p *Params) parse() {
	flag.IntVar(&p.flagPort, "p", 9651, "port to serve")
	flag.StringVar(&p.flagConfigPath, "c", "", "configuration file path")

	flag.Parse()
}

func (e *Engine) setupLogger() {
	e.logger = log.New(os.Stdout, "", log.LUTC|log.Lshortfile|log.Ldate|log.Ltime)
}

func (e *Engine) loadConfig() error {
	var configuration *config.Config

	if len(e.flagConfigPath) > 0 {
		bytes, err := os.ReadFile(e.flagConfigPath)
		if err != nil {
			return err
		}
		configuration, err = config.Decode(bytes)
		if err != nil {
			return err
		}
	} else {
		configuration = config.Default()
	}

	e.config = configuration

	return nil
}

func (e *Engine) serve(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", e.flagPort))
	if err != nil {
		return err
	}
	e.logger.Printf("info: listening %d", e.flagPort)

	var wg sync.WaitGroup

	connC := make(chan net.Conn)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if conn, err := listener.Accept(); err == nil {
					connC <- conn
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			e.logger.Printf("info: all clients done")
			listener.Close()
			return nil
		case conn := <-connC:
			go func() {
				client := client.NewClient(e.logger, conn)

				e.clients.Store(client, struct{}{})

				wg.Add(1)
				client.Serve(ctx)
				wg.Done()

				e.clients.Delete(client)
			}()
		}
	}
}

func main() {
	var engine Engine

	engine.parse()

	engine.setupLogger()

	if err := engine.loadConfig(); err != nil {
		engine.logger.Fatalf("error: %s", err)
	}

	engine.logger.Printf("info: config %s", engine.config)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := engine.serve(ctx); err != nil {
		engine.logger.Fatalf("error: %s", err)
	}
}
