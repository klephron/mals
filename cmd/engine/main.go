package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mals-engine/internal/client"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Params struct {
	flagPort int
}

type Engine struct {
	Params

	logger  *log.Logger
	clients sync.Map
}

func (p *Params) Parse() {
	flag.IntVar(&p.flagPort, "p", 9651, "port to serve")

	flag.Parse()
}

func (e *Engine) SetupLogger() {
	e.logger = log.New(os.Stdout, "", log.LUTC|log.Lshortfile|log.Ldate|log.Ltime)
}

func (e *Engine) Serve(ctx context.Context) error {
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

	engine.Parse()
	engine.SetupLogger()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := engine.Serve(ctx); err != nil {
		engine.logger.Fatal(err)
	}
}
