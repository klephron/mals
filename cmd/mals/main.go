package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mals/internal/client"
	"mals/internal/model"
	"mals/pkg/config"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/puzpuzpuz/xsync/v4"
)

type Params struct {
	flagPort       int
	flagConfigPath string
}

type Engine struct {
	Params
	logger  *log.Logger
	config  *config.Config
	clients *xsync.Map[*client.Client, struct{}]
	models  *xsync.Map[string, model.ModelService] // [id, model]
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

func (e *Engine) setupModels() error {
	e.models = xsync.NewMap[string, model.ModelService]()

	for _, m := range e.config.Models {
		if _, present := e.models.Load(m.Id); present {
			return fmt.Errorf("error: model %s: duplicate id", m.Id)
		}

		if m.Spec != "OpenAI" {
			return fmt.Errorf("error: model %s: spec %s is unsupported", m.Id, m.Spec)
		}

		e.models.Store(m.Id, model.NewModelOpenAI(e.logger, m.Id, m.Spec, m.BaseUrl, m.Settings))
	}
	return nil
}

func (e *Engine) serveModels(ctx context.Context) {
	var wg sync.WaitGroup

	e.models.Range(func(id string, m model.ModelService) bool {
		wg.Add(1)
		go func(m model.ModelService) {
			defer wg.Done()
			m.Serve(ctx)
		}(m)

		return true
	})

	wg.Wait()
	e.logger.Printf("info: all models stopped")
}

func (e *Engine) serveClients(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", e.flagPort))
	if err != nil {
		return err
	}
	e.logger.Printf("info: listening %d started", e.flagPort)

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

	e.clients = xsync.NewMap[*client.Client, struct{}]()

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			e.logger.Printf("info: all clients done")
			listener.Close()
			e.logger.Printf("info: listening stopped")
			return nil
		case conn := <-connC:
			wg.Add(1)
			go func(conn net.Conn) {
				defer wg.Done()

				config, err := e.setupClientConfig(conn)
				if err != nil {
					e.logger.Printf("%s", err)
					return
				}

				client := client.NewClient(e.logger, conn, config)
				e.clients.Store(client, struct{}{})

				client.Serve(ctx)

				e.clients.Delete(client)
			}(conn)
		}
	}
}

func (e *Engine) setupClientConfig(conn net.Conn) (client.Config, error) {
	defaultModelId := e.config.Workspaces.Default.Model.Id
	defaultModel, ok := e.models.Load(defaultModelId)
	if !ok {
		return client.Config{}, fmt.Errorf("error: for client %s model %s not found", conn.LocalAddr(), defaultModelId)
	}
	return client.Config{
		Workspace: client.WorkspaceConfig{
			DefaultModel: defaultModel,
		},
	}, nil
}

func (e *Engine) serve(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		e.serveModels(ctx)
	}()

	go func() {
		defer wg.Done()
		err := e.serveClients(ctx)
		if err != nil {
			e.logger.Printf("error: %s", err)
		}
	}()

	wg.Wait()
}

func main() {
	var engine Engine

	engine.parse()

	engine.setupLogger()

	if err := engine.loadConfig(); err != nil {
		engine.logger.Fatalf("error: %s", err)
	}
	engine.logger.Printf("info: config %s", engine.config)

	if err := engine.setupModels(); err != nil {
		engine.logger.Fatalf("error: %s", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	engine.serve(ctx)
}
